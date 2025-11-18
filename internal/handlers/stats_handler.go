package handlers

import (
	"net/http"
	"reviewer-service/internal/services"

	"github.com/gin-gonic/gin"
)

type StatsHandler struct {
	StatsService services.StatsServiceInterface
}

func NewStatsHandler(statsService services.StatsServiceInterface) *StatsHandler {
	return &StatsHandler{StatsService: statsService}
}

func (h *StatsHandler) GetStats(c *gin.Context) {
	ctx := c.Request.Context()

	stats, err := h.StatsService.GetStats(ctx)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, stats)
}
