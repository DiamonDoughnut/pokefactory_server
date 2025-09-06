package api

import (
	"net/http"
	"time"

	"pokefactory_server/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func (s *Server) healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
		"time":   time.Now(),
	})
}

func (s *Server) login(c *gin.Context) {
	var loginReq struct {
		UUID     string `json:"uuid" binding:"required"`
		Username string `json:"username" binding:"required"`
	}

	if err := c.ShouldBindJSON(&loginReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get or create player
	player, err := s.getOrCreatePlayer(loginReq.UUID, loginReq.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to authenticate player"})
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uuid":      player.UUID,
		"player_id": player.ID,
		"username":  player.Username,
		"exp":       time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.config.JWT.Secret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":  tokenString,
		"player": player,
	})
}

func (s *Server) getPlayerProfile(c *gin.Context) {
	playerUUID := c.GetString("player_uuid")
	
	player, err := s.getPlayerByUUID(playerUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Player not found"})
		return
	}

	c.JSON(http.StatusOK, player)
}

func (s *Server) updatePlayerProfile(c *gin.Context) {
	playerID := c.GetFloat64("player_id")
	
	var updateReq struct {
		Username string `json:"username"`
	}

	if err := c.ShouldBindJSON(&updateReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.updatePlayer(int(playerID), updateReq.Username); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update player"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Player updated successfully"})
}

func (s *Server) getPlayerStats(c *gin.Context) {
	playerID := c.GetFloat64("player_id")
	
	stats, err := s.getPlayerStatsByID(int(playerID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Player stats not found"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (s *Server) updatePlayerStats(c *gin.Context) {
	playerID := c.GetFloat64("player_id")
	
	var stats models.PlayerStats
	if err := c.ShouldBindJSON(&stats); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stats.PlayerID = int(playerID)
	if err := s.updatePlayerStatsByID(stats); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update stats"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Stats updated successfully"})
}

func (s *Server) getPlayerData(c *gin.Context) {
	playerID := c.GetFloat64("player_id")
	key := c.Param("key")
	
	data, err := s.getPlayerDataByKey(int(playerID), key)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"value": data.DataValue})
}

func (s *Server) setPlayerData(c *gin.Context) {
	playerID := c.GetFloat64("player_id")
	key := c.Param("key")
	
	var dataReq struct {
		Value string `json:"value" binding:"required"`
	}

	if err := c.ShouldBindJSON(&dataReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.setPlayerDataByKey(int(playerID), key, dataReq.Value); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data set successfully"})
}

func (s *Server) getPokedexSummary(c *gin.Context) {
	playerID := c.GetFloat64("player_id")
	
	summary, err := s.getPokedexSummaryByID(int(playerID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pokédex summary not found"})
		return
	}

	c.JSON(http.StatusOK, summary)
}

func (s *Server) getRegionalPokedex(c *gin.Context) {
	playerID := c.GetFloat64("player_id")
	
	var regionReq struct {
		Region string `json:"region" binding:"required"`
	}

	if err := c.ShouldBindJSON(&regionReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pokedex, err := s.getRegionalPokedexByID(int(playerID), regionReq.Region)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Regional Pokédex not found"})
		return
	}

	c.JSON(http.StatusOK, pokedex)
}

func (s *Server) updatePokedex(c *gin.Context) {
	playerID := c.GetFloat64("player_id")
	
	var updateReq models.PokedexUpdateRequest
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := s.updatePokedexEntry(int(playerID), updateReq); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update Pokédex"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pokédex updated successfully"})
}

func (s *Server) getPokedexLeaderboard(c *gin.Context) {
	leaderboard, err := s.getPokedexLeaderboardData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get leaderboard"})
		return
	}

	c.JSON(http.StatusOK, leaderboard)
}