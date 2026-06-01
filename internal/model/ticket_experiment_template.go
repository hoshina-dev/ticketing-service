package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TicketExperimentTemplate struct {
	ID                   uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	TicketID             uuid.UUID `gorm:"type:uuid;not null;index:idx_tet_ticket"`
	ExperimentTemplateID uuid.UUID `gorm:"type:uuid;not null;index:idx_tet_experiment_template"`

	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (t *TicketExperimentTemplate) TableName() string {
	return "ticket_experiment_templates"
}
