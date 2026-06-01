package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type HealthHandler struct {
	db *gorm.DB
}

func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

// Live godoc
//
//	@Summary	Liveness probe
//	@Tags		health
//	@Produce	json
//	@Success	200	{object}	map[string]string
//	@Router		/health [get]
func (h *HealthHandler) Live(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"status": "ok"})
}

// Ready godoc
//
//	@Summary	Readiness probe
//	@Tags		health
//	@Produce	json
//	@Success	200	{object}	map[string]string
//	@Failure	503	{object}	map[string]string
//	@Router		/ready [get]
func (h *HealthHandler) Ready(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(c.Context(), 2*time.Second)
	defer cancel()

	sqlDB, err := h.db.DB()
	if err != nil || sqlDB.PingContext(ctx) != nil {
		return c.Status(http.StatusServiceUnavailable).JSON(fiber.Map{"status": "unavailable"})
	}
	return c.JSON(fiber.Map{"status": "ok"})
}
