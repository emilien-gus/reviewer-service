package handlers

import (
	"errors"
	"net/http"
	"reviewer-service/internal/models"
	"reviewer-service/internal/repository"
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
	var errResponce models.ErrorResponse
	if err := c.ShouldBindJSON(&req); err != nil {
		errResponce = models.NewErrorResponse(models.CodeInvalidRequest, err.Error())
		c.JSON(http.StatusBadRequest, errResponce)
		return
	}

	team, err := h.TeamService.CreateTeam(c.Request.Context(), &req)
	if err != nil {
		if errors.Is(repository.ErrTeamExists, err) {
			errResponce = models.NewErrorResponse(models.CodeTeamExists, team.Name+" already exists")
			c.JSON(http.StatusBadRequest, errResponce)
			return
		}

		errResponce = models.NewErrorResponse(models.CodeInternalServerError, err.Error())
		c.JSON(http.StatusInternalServerError, errResponce)
		return
	}

	c.JSON(http.StatusCreated, team)
}

func (h *TeamHandler) GetTeam(c *gin.Context) {
	var errResponse models.ErrorResponse

	teamName := c.Param("team_name")
	if teamName == "" {
		errResponse = models.NewErrorResponse(models.CodeInvalidRequest, "team_name is required")
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	team, err := h.TeamService.GetTeam(c.Request.Context(), teamName)
	if err != nil {
		if errors.Is(err, repository.ErrTeamNotFound) {
			errResponse = models.NewErrorResponse(models.CodeNotFound, err.Error())
			c.JSON(http.StatusNotFound, errResponse)
			return
		}

		errResponse = models.NewErrorResponse(models.CodeInternalServerError, err.Error())
		c.JSON(http.StatusInternalServerError, errResponse)
		return
	}

	c.JSON(http.StatusOK, team)
}
