package handlers

import (
	"errors"
	"net/http"
	"reviewer-service/internal/models"
	"reviewer-service/internal/repository"
	"reviewer-service/internal/services"

	"github.com/gin-gonic/gin"
)

type PullRequestHandler struct {
	service services.PullRequestServiceInterface
}

func NewPullRequestHandler(s services.PullRequestServiceInterface) *PullRequestHandler {
	return &PullRequestHandler{service: s}
}

type createPRRequest struct {
	ID       string `json:"pull_request_id" binding:"required"`
	Name     string `json:"pull_request_name" binding:"required"`
	AuthorID string `json:"author_id" binding:"required"`
}

func (h *PullRequestHandler) CreatePR(c *gin.Context) {
	var req createPRRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		errorResp := models.NewErrorResponse(models.CodeInvalidRequest, err.Error())
		c.JSON(http.StatusBadRequest, errorResp)
		return
	}

	pr, err := h.service.CreatePullRequest(c.Request.Context(), req.ID, req.Name, req.AuthorID)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrPullRequestExists):
			errorResp := models.NewErrorResponse(models.CodePRExists, err.Error())
			c.JSON(http.StatusConflict, errorResp)
		case errors.Is(err, repository.ErrTeamNotFound) || errors.Is(err, repository.ErrUserNotFound):
			errorResp := models.NewErrorResponse(models.CodeNotFound, err.Error())
			c.JSON(http.StatusNotFound, errorResp)
		default:
			errorResp := models.NewErrorResponse(models.CodeInternalServerError, "Internal server error")
			c.JSON(http.StatusInternalServerError, errorResp)
		}
		return
	}

	c.JSON(http.StatusCreated, pr)
}

func (h *PullRequestHandler) SetMergedInPR(c *gin.Context) {
	var req struct {
		ID string `json:"pull_request_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		errorResp := models.NewErrorResponse(models.CodeInvalidRequest, err.Error())
		c.JSON(http.StatusBadRequest, errorResp)
		return
	}

	pr, err := h.service.SetMergedStatus(c.Request.Context(), req.ID)
	if err != nil {
		if errors.Is(err, repository.ErrPullRequestNotFound) {
			errorResp := models.NewErrorResponse(models.CodeNotFound, err.Error())
			c.JSON(http.StatusNotFound, errorResp)
		}

		errorResp := models.NewErrorResponse(models.CodeInternalServerError, err.Error())
		c.JSON(http.StatusInternalServerError, errorResp)
	}

	c.JSON(http.StatusOK, pr)
}
