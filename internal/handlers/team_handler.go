package handlers

import (
	"net/http"
	"reviewer-service/internal/models"
	"reviewer-service/internal/services"

	"github.com/gin-gonic/gin"
)

type TeamHandler struct {
	TeamService services.TeamServiceInterface
}

func NewTeamHandler(TeamService services.TeamServiceInterface) *TeamHandler {
	return &TeamHandler{TeamService: TeamService}
}

func (h *TeamHandler) CreateTeam(c *gin.Context) {
	var req models.Team
	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorReponse(c, http.StatusBadRequest, models.CodeInvalidRequest, err.Error())
		return
	}

	team, err := h.TeamService.CreateTeam(c.Request.Context(), &req)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, team)
}

func (h *TeamHandler) GetTeam(c *gin.Context) {
	teamName := c.Query("team_name")
	if teamName == "" {
		handleErrorReponse(c, http.StatusBadRequest, models.CodeInvalidRequest, "team_name is required")
		return
	}

	team, err := h.TeamService.GetTeam(c.Request.Context(), teamName)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, team)
}
