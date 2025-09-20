package handlers

import (
	"net/http"
	"task-manager/backend/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

type RefreshHandler struct {
	db          *gorm.DB
	authService services.AuthService
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type RefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

func NewRefreshHandler(db *gorm.DB, authService services.AuthService) *RefreshHandler {
	return &RefreshHandler{db: db, authService: authService}
}

func (h *RefreshHandler) Refresh(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse refresh token UUID
	refreshTokenUUID, err := uuid.FromString(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid refresh token format"})
		return
	}

	// Validate refresh token
	token, err := h.authService.ValidateRefreshToken(h.db, refreshTokenUUID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		return
	}

	// Invalidate old refresh token
	if err := h.authService.InvalidateRefreshToken(h.db, refreshTokenUUID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to invalidate refresh token"})
		return
	}

	// Generate new tokens
	accessToken, newRefreshToken, err := h.authService.GenerateToken(h.db, token.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	response := RefreshResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		ExpiresIn:    3600, // 1 hour
	}

	c.JSON(http.StatusOK, response)
}
