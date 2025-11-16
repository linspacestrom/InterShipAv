package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/linspacestrom/InterShipAv/internal/dto"
	"github.com/linspacestrom/InterShipAv/internal/mapper"
	"github.com/linspacestrom/InterShipAv/internal/services"
	"go.uber.org/zap"
)

type UserHandler struct {
	svc  services.UserSer
	logg *zap.Logger
}

func NewUserHandlerStruct(svc services.UserSer, logg *zap.Logger) *UserHandler {
	return &UserHandler{svc: svc, logg: logg}
}

func (h *UserHandler) SetActive(c *gin.Context) {
	var userReq dto.UserRequest
	if err := c.ShouldBindJSON(&userReq); err != nil {
		h.logg.Warn("invalid request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	userDomain := mapper.DTOToUserShort(userReq)

	updatedUser, err := h.svc.SetActive(c.Request.Context(), userDomain.Id, userDomain.IsActive)
	if err != nil {
		h.logg.Error("failed to set active user", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user": mapper.UserToDTO(updatedUser)})
}

func (h *UserHandler) GetReview(c *gin.Context) {
	userId := c.Param("user_id")

	userReview, err := h.svc.GetReview(c.Request.Context(), userId)
	if err != nil {
		h.logg.Error("user not found", zap.Error(err))
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, mapper.DomainReviewToDTOReview(userReview))
}
