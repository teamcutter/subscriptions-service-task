package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/teamcutter/subscriptions-service-task/internal/model"
	"github.com/teamcutter/subscriptions-service-task/internal/repo"
)

type Handler struct {
	repository *repo.SubscriptionRepo
	logger *slog.Logger
}

func NewHandler(repository *repo.SubscriptionRepo, logger *slog.Logger) *Handler {
	return &Handler{
		repository: repository,
		logger: logger,
	}
}

// Create godoc
// @Summary Create new subscription
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body model.Subscription true "Subscription"
// @Success 201
// @Failure 400
// @Router /subscriptions [post]
func (h *Handler) Create(c echo.Context) error {
	var sub model.Subscription
	if err := c.Bind(&sub); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if err := h.repository.Create(&sub); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	h.logger.Info("subscription created",
		"user_id", sub.UserID,
		"service", sub.ServiceName,
		"start", sub.StartDate,
	)
	
	return c.NoContent(http.StatusCreated)
}

// GetAll godoc
// @Summary Get all subscriptions
// @Tags subscriptions
// @Produce json
// @Success 200 {array} model.Subscription
// @Router /subscriptions [get]
func (h *Handler) GetAll(c echo.Context) error {
	subs, err := h.repository.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, subs)
}

// Delete godoc
// @Summary Delete subscription by ID
// @Tags subscriptions
// @Param id path int true "ID"
// @Success 204
// @Router /subscriptions/{id} [delete]
func (h *Handler) Delete(c echo.Context) error {
	idStr := c.Param("id")
	id, _ := strconv.Atoi(idStr)
	if err := h.repository.Delete(id); err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	h.logger.Info("subscription deleted",
		"service_id", id)
	return c.NoContent(http.StatusNoContent)
}

// TotalCost godoc
// @Summary Total of all subscriptions
// @Tags subscriptions
// @Produce json
// @Param user query string true "User ID"
// @Param service query string false "Service name"
// @Param start query string true "Start date (MM-YYYY)"
// @Param end query string true "End date (MM-YYYY)"
// @Success 200 {object} map[string]int
// @Router /subscriptions/total [get]
func (h *Handler) TotalCost(c echo.Context) error {
	userID := c.QueryParam("user")
	service := c.QueryParam("service")
	start := c.QueryParam("start")
	end := c.QueryParam("end")

	total, err := h.repository.TotalCost(userID, service, start, end)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	h.logger.Info("total cost for subscription",
		"user_id", userID,
		"service", service,
		"total", total,
	)
	return c.JSON(http.StatusOK, map[string]int{"total": total})
}
