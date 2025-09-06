package api

import (
	"net/http"
	"strconv"

	"pokefactory_server/internal/models"

	"github.com/gin-gonic/gin"
)

func (s *Server) getWebLeaderboards(c *gin.Context) {
	leaderboards, err := s.getWebLeaderboardData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get leaderboards"})
		return
	}

	c.JSON(http.StatusOK, leaderboards)
}

func (s *Server) getWebPlayerStats(c *gin.Context) {
	username := c.Param("username")
	
	player, err := s.getPlayerByUsername(username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Player not found"})
		return
	}

	stats, err := s.getPlayerStatsByID(player.ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Player stats not found"})
		return
	}

	pokedex, err := s.getPokedexSummaryByID(player.ID)
	if err != nil {
		// Create empty pokedex if not found
		pokedex = &models.PokedexSummary{PlayerID: player.ID}
	}

	// Remove UUID for privacy - only expose username
	publicPlayer := models.Player{
		ID:        player.ID,
		Username:  player.Username,
		LastLogin: player.LastLogin,
		CreatedAt: player.CreatedAt,
		UpdatedAt: player.UpdatedAt,
		// UUID omitted for privacy
	}
	
	webStats := models.WebPlayerStats{
		Player:  publicPlayer,
		Stats:   *stats,
		Pokedex: *pokedex,
	}

	c.JSON(http.StatusOK, webStats)
}

func (s *Server) getWebServerAnalytics(c *gin.Context) {
	analytics, err := s.getServerAnalyticsData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get server analytics"})
		return
	}

	c.JSON(http.StatusOK, analytics)
}

func (s *Server) getWebPokemonPopularity(c *gin.Context) {
	dexParam := c.Param("dex")
	nationalID, err := strconv.Atoi(dexParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid national dex number"})
		return
	}

	popularity, err := s.getPokemonPopularityData(nationalID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pokemon popularity data not found"})
		return
	}

	c.JSON(http.StatusOK, popularity)
}