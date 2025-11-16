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

type PullRequestHandler struct {
	svc  services.PRSer
	logg *zap.Logger
}

func NewPullRequestHandlerStruct(svc services.PRSer, logg *zap.Logger) *PullRequestHandler {
	return &PullRequestHandler{svc: svc, logg: logg}
}

func (h *PullRequestHandler) CreatePR(c *gin.Context) {
	var prDTO dto.PRCreateRequest
	if err := c.ShouldBindJSON(&prDTO); err != nil {
		h.logg.Warn("invalid request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	prDomain := mapper.DTOToPrCreate(prDTO)

	createdPR, err := h.svc.Create(c.Request.Context(), prDomain)

	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, mapper.PRReadToDTO(createdPR))

}

func (h *PullRequestHandler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, validateError.ErrPRExist):
		h.logg.Error("PR already exists", zap.Error(err))
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case errors.Is(err, validateError.UserNotFound):
		h.logg.Error("User not found", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	default:
		h.logg.Error("Internal server error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
