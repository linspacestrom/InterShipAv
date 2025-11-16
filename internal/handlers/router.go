package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/linspacestrom/InterShipAv/internal/services"
	"go.uber.org/zap"
)

func NewTeamHandler(router *gin.Engine, svc services.TeamSer, logg *zap.Logger) {
	h := NewTeamHandlerStruct(svc, logg)

	api := router.Group("/team")
	{
		api.POST("/add", h.CreateTeam)
		api.GET("/get/:team_name", h.GetTeamByName)
	}
}

func NewUserHandler(router *gin.Engine, svc services.UserSer, logg *zap.Logger) {
	h := NewUserHandlerStruct(svc, logg)

	api := router.Group("/users")
	{
		api.POST("/setIsActive", h.SetActive)
		api.GET("/getReview/:user_id", h.GetReview)
	}
}

func NewPullRequestHandler(router *gin.Engine, svc services.PRSer, logg *zap.Logger) {
	h := NewPullRequestHandlerStruct(svc, logg)

	api := router.Group("/pullRequest")
	{
		api.POST("/create", h.CreatePR)
	}
}
