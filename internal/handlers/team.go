package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/linspacestrom/InterShipAv/internal/dto"
	"github.com/linspacestrom/InterShipAv/internal/mapper"
	"github.com/linspacestrom/InterShipAv/internal/services"
	"github.com/linspacestrom/InterShipAv/internal/validateError"
	"go.uber.org/zap"
)

type TeamHandler struct {
	svc  services.TeamSer
	logg *zap.Logger
}

func NewTeamHandlerStruct(svc services.TeamSer, logg *zap.Logger) *TeamHandler {
	return &TeamHandler{svc: svc, logg: logg}
}

func (h *TeamHandler) CreateTeam(c *gin.Context) {
	var teamDTO dto.CreateTeamRequest
	if err := c.ShouldBindJSON(&teamDTO); err != nil {
		h.logg.Warn("invalid request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	teamDomain := mapper.DTOToTeam(teamDTO)

	createdTeam, err := h.svc.Create(c.Request.Context(), teamDomain)
	if err != nil {
		h.logg.Error("failed to create team", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, mapper.CreateTeamToDTO(createdTeam))
}

func (h *TeamHandler) GetTeamByName(c *gin.Context) {
	teamName := c.Param("team_name")

	teamDomain, err := h.svc.GetByName(c.Request.Context(), teamName)
	if err != nil && errors.Is(err, validateError.TeamNotFound) {
		h.logg.Error("not found team", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, mapper.TeamToDTO(teamDomain))
}
