package dto

import (
	"time"

	"github.com/google/uuid"
)

type AddExperimentTemplateRequest struct {
	ExperimentTemplateID uuid.UUID `json:"experiment_template_id" validate:"required"`
}

type TicketExperimentTemplateResponse struct {
	ID                   uuid.UUID `json:"id"`
	TicketID             uuid.UUID `json:"ticket_id"`
	ExperimentTemplateID uuid.UUID `json:"experiment_template_id"`
	CreatedAt            time.Time `json:"created_at"`
}
