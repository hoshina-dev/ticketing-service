package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	_ "github.com/hoshina-dev/ticketing-service/docs"
	"github.com/hoshina-dev/ticketing-service/internal/config"
	"github.com/hoshina-dev/ticketing-service/internal/handler"
	"github.com/hoshina-dev/ticketing-service/internal/middleware"
	"github.com/hoshina-dev/ticketing-service/internal/repository"
	"github.com/hoshina-dev/ticketing-service/internal/router"
	"github.com/hoshina-dev/ticketing-service/internal/service"
)

// @title			Ticketing Service API
// @version		1.0
// @description	Laboratory experiment ticketing service.
// @host			localhost:8080
// @BasePath		/
func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	db, err := repository.NewDB(cfg)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}

	ticketRepo := repository.NewTicketRepository(db)
	tetRepo := repository.NewTicketExperimentTemplateRepository(db)

	ticketSvc := service.NewTicketService(ticketRepo, tetRepo)

	ticketHandler := handler.NewTicketHandler(ticketSvc)
	tetHandler := handler.NewTicketExperimentTemplateHandler(ticketSvc)
	healthHandler := handler.NewHealthHandler(db)

	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})
	app.Use(recover.New())

	router.Register(app, ticketHandler, tetHandler, healthHandler)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		slog.Info("server starting", "port", cfg.Port)
		if err := app.Listen(":" + cfg.Port); err != nil {
			slog.Error("server error", "error", err)
		}
	}()

	<-quit
	slog.Info("shutting down server")
	if err := app.Shutdown(); err != nil {
		slog.Error("shutdown error", "error", err)
	}

	sqlDB, _ := db.DB()
	_ = sqlDB.Close()
	slog.Info("server stopped")
}
