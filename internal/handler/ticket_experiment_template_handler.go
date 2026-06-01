package handler

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/hoshina-dev/ticketing-service/internal/dto"
	"github.com/hoshina-dev/ticketing-service/internal/service"
)

type TicketExperimentTemplateHandler struct {
	svc      service.TicketService
	validate *validator.Validate
}

func NewTicketExperimentTemplateHandler(svc service.TicketService) *TicketExperimentTemplateHandler {
	return &TicketExperimentTemplateHandler{svc: svc, validate: validator.New()}
}

// AddExperimentTemplate godoc
//
//	@Summary	Add an experiment template to a ticket
//	@Tags		experiment-templates
//	@Accept		json
//	@Produce	json
//	@Param		id		path		string								true	"Ticket ID"
//	@Param		body	body		dto.AddExperimentTemplateRequest	true	"Experiment template"
//	@Success	201		{object}	dto.TicketExperimentTemplateResponse
//	@Failure	400		{object}	dto.ErrorResponse
//	@Failure	404		{object}	dto.ErrorResponse
//	@Failure	409		{object}	dto.ErrorResponse
//	@Failure	500		{object}	dto.ErrorResponse
//	@Router		/api/v1/tickets/{id}/experiment-templates [post]
func (h *TicketExperimentTemplateHandler) Add(c *fiber.Ctx) error {
	ticketID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid ticket id")
	}

	var req dto.AddExperimentTemplateRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid request body")
	}
	if err := h.validate.Struct(req); err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	resp, err := h.svc.AddExperimentTemplate(c.Context(), ticketID, req)
	if err != nil {
		return err
	}
	return c.Status(http.StatusCreated).JSON(resp)
}

// RemoveExperimentTemplate godoc
//
//	@Summary	Remove an experiment template from a ticket
//	@Tags		experiment-templates
//	@Produce	json
//	@Param		id			path	string	true	"Ticket ID"
//	@Param		templateId	path	string	true	"Experiment template ID"
//	@Success	204
//	@Failure	400	{object}	dto.ErrorResponse
//	@Failure	404	{object}	dto.ErrorResponse
//	@Failure	500	{object}	dto.ErrorResponse
//	@Router		/api/v1/tickets/{id}/experiment-templates/{templateId} [delete]
func (h *TicketExperimentTemplateHandler) Remove(c *fiber.Ctx) error {
	ticketID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid ticket id")
	}
	templateID, err := uuid.Parse(c.Params("templateId"))
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid template id")
	}

	if err := h.svc.RemoveExperimentTemplate(c.Context(), ticketID, templateID); err != nil {
		return err
	}
	return c.SendStatus(http.StatusNoContent)
}

// ListExperimentTemplates godoc
//
//	@Summary	List experiment templates on a ticket
//	@Tags		experiment-templates
//	@Produce	json
//	@Param		id	path		string	true	"Ticket ID"
//	@Success	200	{array}		dto.TicketExperimentTemplateResponse
//	@Failure	404	{object}	dto.ErrorResponse
//	@Failure	500	{object}	dto.ErrorResponse
//	@Router		/api/v1/tickets/{id}/experiment-templates [get]
func (h *TicketExperimentTemplateHandler) List(c *fiber.Ctx) error {
	ticketID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid ticket id")
	}

	resp, err := h.svc.ListExperimentTemplates(c.Context(), ticketID)
	if err != nil {
		return err
	}
	return c.JSON(resp)
}
