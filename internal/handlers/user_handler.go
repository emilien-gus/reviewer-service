package handlers

import (
	"net/http"
	"reviewer-service/internal/models"
	"reviewer-service/internal/services"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserService services.UserServiceInterface
}

func NewUserHandler(s services.UserServiceInterface) *UserHandler {
	return &UserHandler{UserService: s}
}

type setActiveRequest struct {
	UserID   string `json:"user_id" binding:"required"`
	IsActive *bool  `json:"is_active" binding:"required"`
}

func (h *UserHandler) SetIsActive(c *gin.Context) {
	var req setActiveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handleErrorReponse(c, http.StatusBadRequest, models.CodeInvalidRequest, err.Error())
		return
	}

	user, err := h.UserService.SetIsActive(c.Request.Context(), req.UserID, *req.IsActive)

	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) GetReviews(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		handleErrorReponse(c, http.StatusBadRequest, models.CodeInvalidRequest, "user_id is required")
		return
	}

	reviews, err := h.UserService.GetReviews(c.Request.Context(), userID)

	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user_id":       userID,
		"pull_requests": reviews,
	})
}
