package handlers

import (
	"net/http"
	"reviewer-service/internal/models"
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
		handleErrorReponse(c, http.StatusBadRequest, models.CodeInvalidRequest, err.Error())
		return
	}

	pr, err := h.service.CreatePullRequest(c.Request.Context(), req.ID, req.Name, req.AuthorID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, pr)
}

func (h *PullRequestHandler) SetMergedInPR(c *gin.Context) {
	var req struct {
		ID string `json:"pull_request_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorReponse(c, http.StatusBadRequest, models.CodeInvalidRequest, err.Error())
		return
	}

	pr, err := h.service.SetMergedStatus(c.Request.Context(), req.ID)
	if err != nil {
		handleError(c, err)
	}

	c.JSON(http.StatusOK, pr)
}

func (h *PullRequestHandler) ReassignReviewer(c *gin.Context) {
	var req struct {
		PrID      string `json:"pull_request_id" binding:"required"`
		OldUserID string `json:"old_user_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorReponse(c, http.StatusBadRequest, models.CodeInvalidRequest, err.Error())
		return
	}

	pr, newReviewer, err := h.service.ReassignReviewer(c.Request.Context(), req.PrID, req.OldUserID)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"pr":          pr,
		"replaced_by": newReviewer,
	})
}
