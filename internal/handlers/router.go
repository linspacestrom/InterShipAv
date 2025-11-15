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
