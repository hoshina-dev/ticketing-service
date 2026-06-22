package dto

import (
	"time"

	"github.com/google/uuid"
)

type CreateTicketRequest struct {
	UserID               uuid.UUID `json:"user_id" validate:"required"`
	OrganizationID       uuid.UUID `json:"organization_id" validate:"required"`
	Name                 *string   `json:"name"`
	ExperimentTemplateID uuid.UUID `json:"experiment_template_id" validate:"required"`
}

type TransitionStatusRequest struct {
	Status       string  `json:"status" validate:"required"`
	ClosedReason *string `json:"closed_reason"`
}

type ListTicketsQuery struct {
	UserID         *uuid.UUID `query:"user_id"`
	OrganizationID *uuid.UUID `query:"organization_id"`
	Status         *string    `query:"status"`
	SortBy         string     `query:"sort_by"`
	SortDir        string     `query:"sort_dir"`
}

type TicketResponse struct {
	ID             uuid.UUID `json:"id"`
	UserID         uuid.UUID `json:"user_id"`
	OrganizationID uuid.UUID `json:"organization_id"`
	Name           string    `json:"name"`
	Status         string    `json:"status"`

	ClosedReason        *string    `json:"closed_reason,omitempty"`
	SampleReceivedAt    *time.Time `json:"sample_received_at,omitempty"`
	ExperimentStartedAt *time.Time `json:"experiment_started_at,omitempty"`
	ResultsSubmittedAt  *time.Time `json:"results_submitted_at,omitempty"`
	ClosedAt            *time.Time `json:"closed_at,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	ExperimentTemplate *TicketExperimentTemplateResponse `json:"experiment_template,omitempty"`
}
