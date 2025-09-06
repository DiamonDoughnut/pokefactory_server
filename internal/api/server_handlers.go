package api

import (
	"net/http"

	"pokefactory_server/internal/models"

	"github.com/gin-gonic/gin"
)

// Server authentication for Minecraft server instances
func (s *Server) serverAuth(c *gin.Context) {
	var authReq models.ServerAuthRequest
	if err := c.ShouldBindJSON(&authReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate server credentials (you can customize this logic)
	if authReq.ServerKey != s.config.JWT.Secret {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid server credentials"})
		return
	}

	// Generate server JWT token
	token, err := s.generateServerToken(authReq.ServerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate server token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// Server proxy endpoints for player operations
func (s *Server) serverGetPlayer(c *gin.Context) {
	var req models.ServerPlayerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	player, err := s.getPlayerByUUID(req.PlayerUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Player not found"})
		return
	}

	c.JSON(http.StatusOK, player)
}

func (s *Server) serverCreateOrUpdatePlayer(c *gin.Context) {
	var req models.ServerPlayerUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	player, err := s.getOrCreatePlayer(req.PlayerUUID, req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create/update player"})
		return
	}

	c.JSON(http.StatusOK, player)
}

func (s *Server) serverGetPlayerStats(c *gin.Context) {
	var req models.ServerPlayerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	player, err := s.getPlayerByUUID(req.PlayerUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Player not found"})
		return
	}

	stats, err := s.getPlayerStatsByID(player.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Player stats not found"})
		return
	}

	c.JSON(http.StatusOK, stats)
}

func (s *Server) serverUpdatePlayerStats(c *gin.Context) {
	var req models.ServerPlayerStatsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	player, err := s.getPlayerByUUID(req.PlayerUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Player not found"})
		return
	}

	req.Stats.PlayerID = player.ID
	if err := s.updatePlayerStatsByID(req.Stats); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update stats"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Stats updated successfully"})
}

func (s *Server) serverGetPlayerData(c *gin.Context) {
	var req models.ServerPlayerDataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	player, err := s.getPlayerByUUID(req.PlayerUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Player not found"})
		return
	}

	data, err := s.getPlayerDataByKey(player.ID, req.DataKey)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"value": data.DataValue})
}

func (s *Server) serverSetPlayerData(c *gin.Context) {
	var req models.ServerPlayerDataRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	player, err := s.getPlayerByUUID(req.PlayerUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Player not found"})
		return
	}

	if err := s.setPlayerDataByKey(player.ID, req.DataKey, req.DataValue); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Data set successfully"})
}

func (s *Server) serverGetPokedexSummary(c *gin.Context) {
	var req models.ServerPlayerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	player, err := s.getPlayerByUUID(req.PlayerUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Player not found"})
		return
	}

	summary, err := s.getPokedexSummaryByID(player.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pokédex summary not found"})
		return
	}

	c.JSON(http.StatusOK, summary)
}

func (s *Server) serverGetRegionalPokedex(c *gin.Context) {
	var req models.ServerPokedexRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	player, err := s.getPlayerByUUID(req.PlayerUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Player not found"})
		return
	}

	pokedex, err := s.getRegionalPokedexByID(player.ID, req.Region)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Regional Pokédex not found"})
		return
	}

	c.JSON(http.StatusOK, pokedex)
}

func (s *Server) serverUpdatePokedex(c *gin.Context) {
	var req models.ServerPokedexUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	player, err := s.getPlayerByUUID(req.PlayerUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Player not found"})
		return
	}

	updateReq := models.PokedexUpdateRequest{
		NationalID: req.NationalID,
		Action:     req.Action,
	}

	if err := s.updatePokedexEntry(player.ID, updateReq); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update Pokédex"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Pokédex updated successfully"})
}