package handlers

import (
	"errors"
	"log"
	"net/http"
	"reviewer-service/internal/models"
	"reviewer-service/internal/repository"

	"github.com/gin-gonic/gin"
)

func handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, repository.ErrTeamExists):
		handleErrorReponse(c, http.StatusBadRequest, models.CodeTeamExists, "team_name already exists")

	case errors.Is(err, repository.ErrNotAssigned):
		handleErrorReponse(c, http.StatusConflict, models.CodeNotAssigned, "reviewer is not assigned to this PR")

	case errors.Is(err, repository.ErrPRMerged):
		handleErrorReponse(c, http.StatusConflict, models.CodePRMerged, "cannot reassign on merged PR")

	case errors.Is(err, repository.ErrNoCandidate):
		handleErrorReponse(c, http.StatusConflict, models.CodeNoCandidate, "no active replacement candidate in team")

	case errors.Is(err, repository.ErrPullRequestExists):
		handleErrorReponse(c, http.StatusConflict, models.CodePRExists, "PR id already exists")

	case errors.Is(err, repository.ErrTeamNotFound),
		errors.Is(err, repository.ErrUserNotFound),
		errors.Is(err, repository.ErrPullRequestNotFound):
		handleErrorReponse(c, http.StatusNotFound, models.CodeNotFound, err.Error())

	default:
		handleErrorReponse(c, http.StatusInternalServerError, models.CodeInternalServerError, err.Error())
	}
}

func handleErrorReponse(c *gin.Context, status int, code, message string) {
	errorResp := models.NewErrorResponse(code, message)
	log.Print("ERROR", errorResp)
	c.JSON(status, errorResp)
}
