package model

type TicketStatus string

const (
	StatusRequested     TicketStatus = "REQUESTED"
	StatusPending       TicketStatus = "PENDING"
	StatusExperimenting TicketStatus = "EXPERIMENTING"
	StatusFinalizing    TicketStatus = "FINALIZING"
	StatusClosed        TicketStatus = "CLOSED"
)

// validTransitions defines the allowed forward transitions for each status.
// Transitioning to CLOSED from any non-CLOSED state is handled separately
// (requires closed_reason) and is included here as well.
var validTransitions = map[TicketStatus]map[TicketStatus]bool{
	StatusRequested: {
		StatusPending: true,
		StatusClosed:  true,
	},
	StatusPending: {
		StatusExperimenting: true,
		StatusClosed:        true,
	},
	StatusExperimenting: {
		StatusFinalizing: true,
		StatusClosed:     true,
	},
	StatusFinalizing: {
		StatusClosed: true,
	},
	StatusClosed: {},
}

func (s TicketStatus) CanTransitionTo(target TicketStatus) bool {
	allowed, ok := validTransitions[s]
	if !ok {
		return false
	}
	return allowed[target]
}

// IsManualClose returns true when a transition to CLOSED requires a reason
// (i.e., skipping the normal FINALIZING→CLOSED path).
func (s TicketStatus) IsManualClose(target TicketStatus) bool {
	return target == StatusClosed && s != StatusFinalizing
}

func (s TicketStatus) IsValid() bool {
	switch s {
	case StatusRequested, StatusPending, StatusExperimenting, StatusFinalizing, StatusClosed:
		return true
	}
	return false
}
