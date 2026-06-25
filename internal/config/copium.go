package config

import (
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/hoshina-dev/ticketing-service/internal/model"
)

// CopiumConfig holds static Copium settings loaded from the environment.
// Each ticket stage maps to one Copium email template_id.
type CopiumConfig struct {
	BaseURL   string
	Templates map[model.TicketStatus]uuid.UUID
}

var copiumTemplateEnvKeys = map[model.TicketStatus]string{
	model.StatusRequested:     "COPIUM_TEMPLATE_REQUESTED",
	model.StatusPending:       "COPIUM_TEMPLATE_PENDING",
	model.StatusExperimenting: "COPIUM_TEMPLATE_EXPERIMENTING",
	model.StatusFinalizing:    "COPIUM_TEMPLATE_FINALIZING",
	model.StatusClosed:        "COPIUM_TEMPLATE_CLOSED",
}

func loadCopiumConfig() (CopiumConfig, error) {
	cfg := CopiumConfig{
		BaseURL:   getEnv("COPIUM_BASE_URL", ""),
		Templates: make(map[model.TicketStatus]uuid.UUID, len(copiumTemplateEnvKeys)),
	}

	for status, envKey := range copiumTemplateEnvKeys {
		raw := os.Getenv(envKey)
		if raw == "" {
			return CopiumConfig{}, fmt.Errorf("%s is required", envKey)
		}

		id, err := uuid.Parse(raw)
		if err != nil {
			return CopiumConfig{}, fmt.Errorf("%s must be a valid UUID: %w", envKey, err)
		}

		cfg.Templates[status] = id
	}

	return cfg, nil
}

// TemplateForStatus returns the Copium template_id configured for a ticket stage.
func (c CopiumConfig) TemplateForStatus(status model.TicketStatus) (uuid.UUID, bool) {
	id, ok := c.Templates[status]
	return id, ok
}
