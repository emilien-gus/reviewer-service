package handlers

import (
	"database/sql"
	"reviewer-service/internal/repository"
	"reviewer-service/internal/services"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(db *sql.DB, r *gin.Engine) {
	// Инициализация репозиториев
	teamRepo := repository.NewTeamRepository(db)
	userRepo := repository.NewUserRepository(db)
	prRepo := repository.NewPullRequestRepository(db)
	statsRepo := repository.NewStatsRepository(db)

	// Инициализация сервисов
	teamService := services.NewTeamService(teamRepo)
	userService := services.NewUserService(userRepo)
	prService := services.NewPullRequestService(prRepo)
	statsSevice := services.NewStatsService(*statsRepo)

	// Инициализация хэндлеров
	teamHandler := NewTeamHandler(teamService)
	userHandler := NewUserHandler(userService)
	prHandler := NewPullRequestHandler(prService)
	StatsHandler := NewStatsHandler(statsSevice)

	// Группа API с аутентификацией
	api := r.Group("/")
	{
		// Teams
		api.POST("/team/add", teamHandler.CreateTeam)
		api.GET("/team/get", teamHandler.GetTeam)

		// Users
		api.POST("/users/setIsActive", userHandler.SetIsActive)
		api.GET("/users/getReview", userHandler.GetReviews)

		// Pull Requests
		api.POST("/pullRequest/create", prHandler.CreatePR)
		api.POST("/pullRequest/merge", prHandler.SetMergedInPR)
		api.POST("/pullRequest/reassign", prHandler.ReassignReviewer)

		// Stats
		api.GET("/stats", StatsHandler.GetStats)
	}
}
