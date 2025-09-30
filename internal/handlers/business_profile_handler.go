package handlers

import (
	"net/http"

	"github.com/codeZe-us/vestroll-backend/internal/models"
	"github.com/codeZe-us/vestroll-backend/internal/services"
	"github.com/gin-gonic/gin"
)

// BusinessProfileHandler manages business profile endpoints

type BusinessProfileHandler struct {
	service *services.BusinessProfileService
}

func NewBusinessProfileHandler(service *services.BusinessProfileService) *BusinessProfileHandler {
	return &BusinessProfileHandler{service: service}
}

// RegisterRoutes registers business profile routes under /api/v1/profile
func (h *BusinessProfileHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/business-details", h.PostBusinessDetails)
}

// PostBusinessDetails handles POST /api/v1/profile/business-details
func (h *BusinessProfileHandler) PostBusinessDetails(c *gin.Context) {
	var req models.BusinessDetailsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	// Only contractor supported for now
	if req.AccountType != models.AccountTypeContractor {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "validation_error",
			Message: "unsupported account_type",
		})
		return
	}

	if err := h.service.ValidateContractor(req); err != nil {
		status := http.StatusBadRequest
		c.JSON(status, models.ErrorResponse{
			Error:   "validation_error",
			Message: err.Error(),
		})
		return
	}

	profile := h.service.BuildProfile(req)

	if err := h.service.Save(c.Request.Context(), profile); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error:   "internal_error",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.BusinessProfileResponse{
		Success: true,
		Message: "Business details saved",
		Profile: profile,
	})
}