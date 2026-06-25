package copium

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/hoshina-dev/ticketing-service/internal/model"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// Notifier enqueues Copium emails when tickets advance to a new stage.
type Notifier struct {
	client    *ClientWithResponses
	templates map[model.TicketStatus]uuid.UUID
}

func NewNotifier(baseURL string, templates map[model.TicketStatus]uuid.UUID) (*Notifier, error) {
	client, err := NewClientWithResponses(baseURL)
	if err != nil {
		return nil, fmt.Errorf("create copium client: %w", err)
	}

	return &Notifier{
		client:    client,
		templates: templates,
	}, nil
}

// NotifyStageAdvanced enqueues an email for the ticket owner via Copium (user_id mode).
func (n *Notifier) NotifyStageAdvanced(ctx context.Context, ticket *model.Ticket, stage model.TicketStatus) error {
	templateID, ok := n.templates[stage]
	if !ok {
		return fmt.Errorf("no copium template configured for stage %s", stage)
	}

	userID := openapi_types.UUID(ticket.UserID)
	params := ticketNotificationParams(ticket)

	resp, err := n.client.PostEmailsSendWithResponse(ctx, SendEmailRequest{
		TemplateId: openapi_types.UUID(templateID),
		UserId:     &userID,
		Params:     &params,
	})
	if err != nil {
		return fmt.Errorf("copium send email: %w", err)
	}

	if resp.StatusCode() == http.StatusAccepted {
		return nil
	}

	return fmt.Errorf("copium send email: %s (status %d)", errorMessage(resp), resp.StatusCode())
}

func errorMessage(resp *PostEmailsSendResponse) string {
	switch {
	case resp.JSON400 != nil && resp.JSON400.Error != nil:
		return *resp.JSON400.Error
	case resp.JSON404 != nil && resp.JSON404.Error != nil:
		return *resp.JSON404.Error
	case resp.JSON502 != nil && resp.JSON502.Error != nil:
		return *resp.JSON502.Error
	default:
		return resp.Status()
	}
}
