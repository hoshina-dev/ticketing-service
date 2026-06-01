package handler

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/hoshina-dev/ticketing-service/internal/dto"
	"github.com/hoshina-dev/ticketing-service/internal/service"
)

type TicketHandler struct {
	svc      service.TicketService
	validate *validator.Validate
}

func NewTicketHandler(svc service.TicketService) *TicketHandler {
	return &TicketHandler{svc: svc, validate: validator.New()}
}

// CreateTicket godoc
//
//	@Summary	Create a ticket
//	@Tags		tickets
//	@Accept		json
//	@Produce	json
//	@Param		body	body		dto.CreateTicketRequest	true	"Create ticket"
//	@Success	201		{object}	dto.TicketResponse
//	@Failure	400		{object}	dto.ErrorResponse
//	@Failure	500		{object}	dto.ErrorResponse
//	@Router		/api/v1/tickets [post]
func (h *TicketHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateTicketRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	if err := h.validate.Struct(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	resp, err := h.svc.CreateTicket(c.Context(), req)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(resp)
}

// ListTickets godoc
//
//	@Summary	List tickets
//	@Tags		tickets
//	@Produce	json
//	@Param		user_id			query		string	false	"Filter by user ID"
//	@Param		organization_id	query		string	false	"Filter by organization ID"
//	@Param		status			query		string	false	"Filter by status"
//	@Param		sort_by			query		string	false	"Sort field (created_at, updated_at, status)"
//	@Param		sort_dir		query		string	false	"Sort direction (asc, desc)"
//	@Success	200				{array}		dto.TicketResponse
//	@Failure	400				{object}	dto.ErrorResponse
//	@Failure	500				{object}	dto.ErrorResponse
//	@Router		/api/v1/tickets [get]
func (h *TicketHandler) List(c *fiber.Ctx) error {
	var q dto.ListTicketsQuery
	if err := c.QueryParser(&q); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid query parameters")
	}

	resp, err := h.svc.ListTickets(c.Context(), q)
	if err != nil {
		return err
	}
	return c.JSON(resp)
}

// GetByID godoc
//
//	@Summary	Get a ticket by ID
//	@Tags		tickets
//	@Produce	json
//	@Param		id	path		string	true	"Ticket ID"
//	@Success	200	{object}	dto.TicketResponse
//	@Failure	404	{object}	dto.ErrorResponse
//	@Failure	500	{object}	dto.ErrorResponse
//	@Router		/api/v1/tickets/{id} [get]
func (h *TicketHandler) GetByID(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid ticket id")
	}

	resp, err := h.svc.GetTicket(c.Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(resp)
}

// TransitionStatus godoc
//
//	@Summary	Transition ticket status
//	@Tags		tickets
//	@Accept		json
//	@Produce	json
//	@Param		id		path		string						true	"Ticket ID"
//	@Param		body	body		dto.TransitionStatusRequest	true	"Transition request"
//	@Success	200		{object}	dto.TicketResponse
//	@Failure	400		{object}	dto.ErrorResponse
//	@Failure	404		{object}	dto.ErrorResponse
//	@Failure	409		{object}	dto.ErrorResponse
//	@Failure	422		{object}	dto.ErrorResponse
//	@Failure	500		{object}	dto.ErrorResponse
//	@Router		/api/v1/tickets/{id}/status [patch]
func (h *TicketHandler) TransitionStatus(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid ticket id")
	}

	var req dto.TransitionStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	if err := h.validate.Struct(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	resp, err := h.svc.TransitionStatus(c.Context(), id, req)
	if err != nil {
		return err
	}
	return c.JSON(resp)
}

// Delete godoc
//
//	@Summary	Delete a ticket
//	@Tags		tickets
//	@Produce	json
//	@Param		id	path	string	true	"Ticket ID"
//	@Success	204
//	@Failure	404	{object}	dto.ErrorResponse
//	@Failure	500	{object}	dto.ErrorResponse
//	@Router		/api/v1/tickets/{id} [delete]
func (h *TicketHandler) Delete(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid ticket id")
	}

	if err := h.svc.DeleteTicket(c.Context(), id); err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}
