package api

import (
	"net/http"

	"pokefactory_server/internal/models"

	"github.com/gin-gonic/gin"
)

// Simplified endpoint for national dex number updates
func (s *Server) updatePokedexSimple(c *gin.Context) {
	playerID := c.GetFloat64("player_id")
	
	var simpleReq struct {
		NationalID int    `json:"national_id" binding:"required"`
		Action     string `json:"action" binding:"required"` // "catch" or "see"
	}

	if err := c.ShouldBindJSON(&simpleReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert to full request
	req := models.PokedexUpdateRequest{
		NationalID: simpleReq.NationalID,
		Action:     simpleReq.Action,
	}

	if err := s.updatePokedexEntry(int(playerID), req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update Pokédex"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pokédex updated successfully"})
}