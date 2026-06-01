package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/hoshina-dev/ticketing-service/internal/apperr"
	"github.com/hoshina-dev/ticketing-service/internal/model"
"gorm.io/gorm"
)

var allowedSortColumns = map[string]string{
	"created_at": "tickets.created_at",
	"updated_at": "tickets.updated_at",
	"status":     "tickets.status",
}

type TicketFilter struct {
	Status         *model.TicketStatus
	OrganizationID *uuid.UUID
	UserID         *uuid.UUID
}

type TicketRepository interface {
	Create(ctx context.Context, ticket *model.Ticket) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Ticket, error)
	GetAll(ctx context.Context, filter TicketFilter, sortBy, sortDir string) ([]*model.Ticket, error)
	Update(ctx context.Context, ticket *model.Ticket) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type ticketRepository struct {
	db *gorm.DB
}

func NewTicketRepository(db *gorm.DB) TicketRepository {
	return &ticketRepository{db: db}
}

func (r *ticketRepository) Create(ctx context.Context, ticket *model.Ticket) error {
	return r.db.WithContext(ctx).Create(ticket).Error
}

func (r *ticketRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.Ticket, error) {
	var ticket model.Ticket
	err := r.db.WithContext(ctx).
		Preload("ExperimentTemplates").
		First(&ticket, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.ErrTicketNotFound
		}
		return nil, err
	}
	return &ticket, nil
}

func (r *ticketRepository) GetAll(ctx context.Context, filter TicketFilter, sortBy, sortDir string) ([]*model.Ticket, error) {
	q := r.db.WithContext(ctx).Model(&model.Ticket{})

	if filter.Status != nil {
		q = q.Where("status = ?", *filter.Status)
	}
	if filter.OrganizationID != nil {
		q = q.Where("organization_id = ?", *filter.OrganizationID)
	}
	if filter.UserID != nil {
		q = q.Where("user_id = ?", *filter.UserID)
	}

	col, ok := allowedSortColumns[sortBy]
	if !ok {
		col = "tickets.created_at"
		sortDir = "desc"
	}
	if sortDir != "asc" && sortDir != "desc" {
		sortDir = "desc"
	}

	var tickets []*model.Ticket
	err := q.
		Preload("ExperimentTemplates").
		Order(fmt.Sprintf("%s %s", col, sortDir)).
		Find(&tickets).Error
	if err != nil {
		return nil, err
	}

	return tickets, nil
}

func (r *ticketRepository) Update(ctx context.Context, ticket *model.Ticket) error {
	return r.db.WithContext(ctx).Save(ticket).Error
}

func (r *ticketRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Delete(&model.Ticket{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return apperr.ErrTicketNotFound
	}
	return nil
}

