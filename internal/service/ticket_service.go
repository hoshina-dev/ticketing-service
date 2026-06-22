package service

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hoshina-dev/ticketing-service/internal/apperr"
	"github.com/hoshina-dev/ticketing-service/internal/dto"
	"github.com/hoshina-dev/ticketing-service/internal/model"
	"github.com/hoshina-dev/ticketing-service/internal/repository"
)

type TicketService interface {
	CreateTicket(ctx context.Context, req dto.CreateTicketRequest) (*dto.TicketResponse, error)
	GetTicket(ctx context.Context, id uuid.UUID) (*dto.TicketResponse, error)
	ListTickets(ctx context.Context, q dto.ListTicketsQuery) ([]*dto.TicketResponse, error)
	TransitionStatus(ctx context.Context, id uuid.UUID, req dto.TransitionStatusRequest) (*dto.TicketResponse, error)
	DeleteTicket(ctx context.Context, id uuid.UUID) error
	AddExperimentTemplate(ctx context.Context, ticketID uuid.UUID, req dto.AddExperimentTemplateRequest) (*dto.TicketExperimentTemplateResponse, error)
	RemoveExperimentTemplate(ctx context.Context, ticketID, experimentTemplateID uuid.UUID) error
	ListExperimentTemplates(ctx context.Context, ticketID uuid.UUID) ([]*dto.TicketExperimentTemplateResponse, error)
}

type ticketService struct {
	ticketRepo repository.TicketRepository
	tetRepo    repository.TicketExperimentTemplateRepository
}

func NewTicketService(
	ticketRepo repository.TicketRepository,
	tetRepo repository.TicketExperimentTemplateRepository,
) TicketService {
	return &ticketService{
		ticketRepo: ticketRepo,
		tetRepo:    tetRepo,
	}
}

func (s *ticketService) CreateTicket(ctx context.Context, req dto.CreateTicketRequest) (*dto.TicketResponse, error) {
	name := generateTicketName()
	if req.Name != nil && strings.TrimSpace(*req.Name) != "" {
		name = strings.TrimSpace(*req.Name)
	}

	ticket := &model.Ticket{
		UserID:         req.UserID,
		OrganizationID: req.OrganizationID,
		Name:           name,
		Status:         model.StatusRequested,
	}

	if err := s.ticketRepo.Create(ctx, ticket); err != nil {
		return nil, err
	}

	tet := &model.TicketExperimentTemplate{
		TicketID:             ticket.ID,
		ExperimentTemplateID: req.ExperimentTemplateID,
	}
	if err := s.tetRepo.Add(ctx, tet); err != nil {
		return nil, err
	}

	created, err := s.ticketRepo.GetByID(ctx, ticket.ID)
	if err != nil {
		return nil, err
	}

	return toTicketResponse(created), nil
}

func (s *ticketService) GetTicket(ctx context.Context, id uuid.UUID) (*dto.TicketResponse, error) {
	ticket, err := s.ticketRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return toTicketResponse(ticket), nil
}

func (s *ticketService) ListTickets(ctx context.Context, q dto.ListTicketsQuery) ([]*dto.TicketResponse, error) {
	filter := repository.TicketFilter{
		UserID:         q.UserID,
		OrganizationID: q.OrganizationID,
	}

	if q.Status != nil {
		status := model.TicketStatus(*q.Status)
		if !status.IsValid() {
			return nil, apperr.ErrInvalidStatus
		}
		filter.Status = &status
	}

	tickets, err := s.ticketRepo.GetAll(ctx, filter, q.SortBy, q.SortDir)
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.TicketResponse, len(tickets))
	for i, t := range tickets {
		responses[i] = toTicketResponse(t)
	}
	return responses, nil
}

func (s *ticketService) TransitionStatus(ctx context.Context, id uuid.UUID, req dto.TransitionStatusRequest) (*dto.TicketResponse, error) {
	ticket, err := s.ticketRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if ticket.Status == model.StatusClosed {
		return nil, apperr.ErrTicketClosed
	}

	target := model.TicketStatus(req.Status)
	if !target.IsValid() {
		return nil, apperr.ErrInvalidStatus
	}

	if !ticket.Status.CanTransitionTo(target) {
		return nil, apperr.ErrInvalidTransition
	}

	if ticket.Status.IsManualClose(target) {
		if req.ClosedReason == nil || *req.ClosedReason == "" {
			return nil, apperr.ErrClosedReasonRequired
		}
		ticket.ClosedReason = req.ClosedReason
	}

	now := time.Now()
	switch {
	case ticket.Status == model.StatusRequested && target == model.StatusPending:
		ticket.SampleReceivedAt = &now
	case ticket.Status == model.StatusPending && target == model.StatusExperimenting:
		ticket.ExperimentStartedAt = &now
	case ticket.Status == model.StatusExperimenting && target == model.StatusFinalizing:
		ticket.ResultsSubmittedAt = &now
	case target == model.StatusClosed:
		ticket.ClosedAt = &now
	}

	ticket.Status = target

	if err := s.ticketRepo.Update(ctx, ticket); err != nil {
		return nil, err
	}

	return toTicketResponse(ticket), nil
}

func (s *ticketService) DeleteTicket(ctx context.Context, id uuid.UUID) error {
	return s.ticketRepo.Delete(ctx, id)
}

func (s *ticketService) AddExperimentTemplate(ctx context.Context, ticketID uuid.UUID, req dto.AddExperimentTemplateRequest) (*dto.TicketExperimentTemplateResponse, error) {
	ticket, err := s.ticketRepo.GetByID(ctx, ticketID)
	if err != nil {
		return nil, err
	}
	if ticket.Status == model.StatusClosed {
		return nil, apperr.ErrTicketClosed
	}

	existing, err := s.tetRepo.GetByTicketID(ctx, ticketID)
	if err != nil {
		return nil, err
	}
	if len(existing) > 0 {
		return nil, apperr.ErrTicketHasTemplate
	}

	tet := &model.TicketExperimentTemplate{
		TicketID:             ticketID,
		ExperimentTemplateID: req.ExperimentTemplateID,
	}
	if err := s.tetRepo.Add(ctx, tet); err != nil {
		return nil, err
	}

	return toTETResponse(tet), nil
}

func (s *ticketService) RemoveExperimentTemplate(ctx context.Context, ticketID, experimentTemplateID uuid.UUID) error {
	ticket, err := s.ticketRepo.GetByID(ctx, ticketID)
	if err != nil {
		return err
	}
	if ticket.Status == model.StatusClosed {
		return apperr.ErrTicketClosed
	}

	return s.tetRepo.Remove(ctx, ticketID, experimentTemplateID)
}

func (s *ticketService) ListExperimentTemplates(ctx context.Context, ticketID uuid.UUID) ([]*dto.TicketExperimentTemplateResponse, error) {
	if _, err := s.ticketRepo.GetByID(ctx, ticketID); err != nil {
		return nil, err
	}

	items, err := s.tetRepo.GetByTicketID(ctx, ticketID)
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.TicketExperimentTemplateResponse, len(items))
	for i, item := range items {
		responses[i] = toTETResponse(item)
	}
	return responses, nil
}

func toTicketResponse(t *model.Ticket) *dto.TicketResponse {
	resp := &dto.TicketResponse{
		ID:                  t.ID,
		UserID:              t.UserID,
		OrganizationID:      t.OrganizationID,
		Name:                t.Name,
		Status:              string(t.Status),
		ClosedReason:        t.ClosedReason,
		SampleReceivedAt:    t.SampleReceivedAt,
		ExperimentStartedAt: t.ExperimentStartedAt,
		ResultsSubmittedAt:  t.ResultsSubmittedAt,
		ClosedAt:            t.ClosedAt,
		CreatedAt:           t.CreatedAt,
		UpdatedAt:           t.UpdatedAt,
	}
	if t.ExperimentTemplate != nil {
		resp.ExperimentTemplate = toTETResponse(t.ExperimentTemplate)
	}
	return resp
}

func toTETResponse(t *model.TicketExperimentTemplate) *dto.TicketExperimentTemplateResponse {
	return &dto.TicketExperimentTemplateResponse{
		ID:                   t.ID,
		TicketID:             t.TicketID,
		ExperimentTemplateID: t.ExperimentTemplateID,
		CreatedAt:            t.CreatedAt,
	}
}
