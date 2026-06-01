package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Ticket struct {
	ID             uuid.UUID    `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID         uuid.UUID    `gorm:"type:uuid;not null;index:idx_tickets_user"`
	OrganizationID uuid.UUID    `gorm:"type:uuid;not null;index:idx_tickets_organization"`
	Status         TicketStatus `gorm:"type:ticket_status;not null;default:REQUESTED;index:idx_tickets_status"`
	ClosedReason   *string      `gorm:"type:text"`

	SampleReceivedAt    *time.Time `gorm:"type:timestamptz"`
	ExperimentStartedAt *time.Time `gorm:"type:timestamptz"`
	ResultsSubmittedAt  *time.Time `gorm:"type:timestamptz"`
	ClosedAt            *time.Time `gorm:"type:timestamptz"`

	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	ExperimentTemplates []TicketExperimentTemplate `gorm:"foreignKey:TicketID"`
}

func (t *Ticket) TableName() string {
	return "tickets"
}
