package copium

import (
	"time"

	"github.com/hoshina-dev/ticketing-service/internal/model"
)

func ticketNotificationParams(ticket *model.Ticket) map[string]interface{} {
	return map[string]interface{}{
		"closed_at":             formatOptionalTime(ticket.ClosedAt),
		"created_at":            formatTime(ticket.CreatedAt),
		"experiment_started_at": formatOptionalTime(ticket.ExperimentStartedAt),
		"result_submitted_at":   formatOptionalTime(ticket.ResultsSubmittedAt),
		"sample_received_at":    formatOptionalTime(ticket.SampleReceivedAt),
		"status":                string(ticket.Status),
		"ticket_name":           ticket.Name,
		"updated_at":            formatTime(ticket.UpdatedAt),
	}
}

func formatOptionalTime(t *time.Time) string {
	if t == nil {
		return "-"
	}
	return t.UTC().Format(time.RFC3339)
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return "-"
	}
	return t.UTC().Format(time.RFC3339)
}
