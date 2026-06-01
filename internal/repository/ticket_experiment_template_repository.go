package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hoshina-dev/ticketing-service/internal/apperr"
	"github.com/hoshina-dev/ticketing-service/internal/model"
	"gorm.io/gorm"
)

type TicketExperimentTemplateRepository interface {
	Add(ctx context.Context, tet *model.TicketExperimentTemplate) error
	Remove(ctx context.Context, ticketID, experimentTemplateID uuid.UUID) error
	GetByTicketID(ctx context.Context, ticketID uuid.UUID) ([]*model.TicketExperimentTemplate, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.TicketExperimentTemplate, error)
}

type tetRepository struct {
	db *gorm.DB
}

func NewTicketExperimentTemplateRepository(db *gorm.DB) TicketExperimentTemplateRepository {
	return &tetRepository{db: db}
}

func (r *tetRepository) Add(ctx context.Context, tet *model.TicketExperimentTemplate) error {
	err := r.db.WithContext(ctx).Create(tet).Error
	if err != nil {
		if isPgUniqueViolation(err) {
			return apperr.ErrDuplicateTemplate
		}
		return err
	}
	return nil
}

func (r *tetRepository) Remove(ctx context.Context, ticketID, experimentTemplateID uuid.UUID) error {
	result := r.db.WithContext(ctx).
		Where("ticket_id = ? AND experiment_template_id = ?", ticketID, experimentTemplateID).
		Delete(&model.TicketExperimentTemplate{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return apperr.ErrTemplateNotFound
	}
	return nil
}

func (r *tetRepository) GetByTicketID(ctx context.Context, ticketID uuid.UUID) ([]*model.TicketExperimentTemplate, error) {
	var items []*model.TicketExperimentTemplate
	err := r.db.WithContext(ctx).
		Where("ticket_id = ?", ticketID).
		Find(&items).Error
	return items, err
}

func (r *tetRepository) GetByID(ctx context.Context, id uuid.UUID) (*model.TicketExperimentTemplate, error) {
	var item model.TicketExperimentTemplate
	err := r.db.WithContext(ctx).First(&item, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.ErrTemplateNotFound
		}
		return nil, err
	}
	return &item, nil
}
