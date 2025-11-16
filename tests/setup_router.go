package tests

import (
	"github.com/gin-gonic/gin"
	"github.com/linspacestrom/InterShipAv/internal/handlers"
	"github.com/linspacestrom/InterShipAv/internal/services"
	"go.uber.org/zap"
)

func SetupTestRouter(team services.TeamSer, user services.UserSer, pr services.PRSer) *gin.Engine {
	r := gin.Default()
	logger := zap.NewNop()

	hTeam := handlers.NewTeamHandlerStruct(team, logger)
	teamAPI := r.Group("/team")
	teamAPI.POST("/add", hTeam.CreateTeam)
	teamAPI.GET("/get", hTeam.GetTeamByName)

	hUser := handlers.NewUserHandlerStruct(user, logger)
	userAPI := r.Group("/users")
	userAPI.POST("/setIsActive", hUser.SetActive)
	userAPI.GET("/getReview", hUser.GetReview)

	hPR := handlers.NewPullRequestHandlerStruct(pr, logger)
	prAPI := r.Group("/pullRequest")
	prAPI.POST("/create", hPR.CreatePR)
	prAPI.POST("/merge", hPR.MergePR)
	prAPI.POST("/reassign", hPR.ReassignPR)

	return r
}
