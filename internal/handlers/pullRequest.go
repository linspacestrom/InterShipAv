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
		h.logg.Error("invalid request", zap.Error(err))
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

func (h *PullRequestHandler) MergePR(c *gin.Context) {
	var prDTO dto.PRMergeRequest
	if err := c.ShouldBindJSON(&prDTO); err != nil {
		h.logg.Error("invalid request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	prDomain := mapper.DTOtoPRMerge(prDTO)

	mergedPr, err := h.svc.Merge(c.Request.Context(), prDomain)
	if err != nil {
		h.handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"pr": mapper.PRMergeToDTO(mergedPr)})
}

func (h *PullRequestHandler) ReassignPR(c *gin.Context) {
	var prDTO dto.PRReassignRequest
	if err := c.ShouldBindJSON(&prDTO); err != nil {
		h.logg.Error("invalid request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	prDomain := mapper.PrReassignDTOtoDomain(prDTO)

	updatedPR, err := h.svc.Reassign(c.Request.Context(), prDomain)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, mapper.DomainToPRDTO(updatedPR))
}

func (h *PullRequestHandler) handleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, validateError.ErrTeamExists):
		h.logg.Error("Team already exists", zap.Error(err))
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})

	case errors.Is(err, validateError.TeamNotFound):
		h.logg.Error("Team not found", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

	case errors.Is(err, validateError.UserNotFound):
		h.logg.Error("User not found", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

	case errors.Is(err, validateError.UserNotAssignToTeam):
		h.logg.Error("User not assigned to team", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

	case errors.Is(err, validateError.ErrPRExist):
		h.logg.Error("Pull request already exists", zap.Error(err))
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})

	case errors.Is(err, validateError.ErrPrNotExist):
		h.logg.Error("Pull request not found", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})

	case errors.Is(err, validateError.PrMergedExist):
		h.logg.Error("Pull request already merged", zap.Error(err))
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})

	case errors.Is(err, validateError.UserNotAssignReviewer):
		h.logg.Error("User not assigned as reviewer", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

	case errors.Is(err, validateError.NoCandidate):
		h.logg.Error("No active replacement candidate", zap.Error(err))
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})

	default:
		h.logg.Error("Internal server error", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
