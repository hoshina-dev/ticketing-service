package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/hoshina-dev/ticketing-service/internal/handler"
	"github.com/hoshina-dev/ticketing-service/internal/middleware"
)

func Register(
	app *fiber.App,
	ticketHandler *handler.TicketHandler,
	tetHandler *handler.TicketExperimentTemplateHandler,
	healthHandler *handler.HealthHandler,
) {
	app.Use(middleware.Logger())

	app.Get("/health", healthHandler.Live)
	app.Get("/ready", healthHandler.Ready)
	app.Get("/docs/*", swagger.HandlerDefault)

	v1 := app.Group("/api/v1")

	tickets := v1.Group("/tickets")
	tickets.Post("/", ticketHandler.Create)
	tickets.Get("/", ticketHandler.List)
	tickets.Get("/:id", ticketHandler.GetByID)
	tickets.Patch("/:id/status", ticketHandler.TransitionStatus)
	tickets.Delete("/:id", ticketHandler.Delete)

	tickets.Get("/:id/experiment-templates", tetHandler.List)
	tickets.Post("/:id/experiment-templates", tetHandler.Add)
	tickets.Delete("/:id/experiment-templates/:templateId", tetHandler.Remove)
}
