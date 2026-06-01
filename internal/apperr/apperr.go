package apperr

import "net/http"

type Error struct {
	HTTPStatus int
	Message    string
}

func (e *Error) Error() string {
	return e.Message
}

var (
	ErrTicketNotFound       = &Error{http.StatusNotFound, "ticket not found"}
	ErrTemplateNotFound     = &Error{http.StatusNotFound, "experiment template not found on this ticket"}
	ErrDuplicateTemplate    = &Error{http.StatusConflict, "experiment template already added to this ticket"}
	ErrTicketClosed         = &Error{http.StatusConflict, "ticket is already closed"}
	ErrInvalidTransition    = &Error{http.StatusUnprocessableEntity, "status transition is not allowed"}
	ErrClosedReasonRequired = &Error{http.StatusUnprocessableEntity, "closed_reason is required when closing a ticket manually"}
	ErrInvalidStatus        = &Error{http.StatusBadRequest, "invalid ticket status value"}
)
