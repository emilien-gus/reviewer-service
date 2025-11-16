package handlers

import (
	"errors"
	"net/http"
	"reviewer-service/internal/models"
	"reviewer-service/internal/repository"
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
	IsActive bool   `json:"is_active" binding:"required"`
}

func (h *UserHandler) SetIsActive(c *gin.Context) {
	var req setActiveRequest
	var errResponse models.ErrorResponse
	if err := c.ShouldBindJSON(&req); err != nil {
		errResponse = models.NewErrorResponse(models.CodeInvalidRequest, err.Error())
		c.JSON(http.StatusBadRequest, errResponse)
		return
	}

	user, err := h.UserService.SetIsActive(
		c.Request.Context(),
		req.UserID,
		req.IsActive,
	)

	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			errResp := models.NewErrorResponse(models.CodeNotFound, err.Error())
			c.JSON(http.StatusNotFound, errResp)
			return
		}

		errResp := models.NewErrorResponse(models.CodeInternalServerError, err.Error())
		c.JSON(http.StatusInternalServerError, errResp)
		return
	}

	c.JSON(http.StatusOK, user)
}
