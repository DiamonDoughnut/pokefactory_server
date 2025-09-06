package api

import (
	"database/sql"

	"pokefactory_server/internal/config"
	"pokefactory_server/internal/middleware"

	"github.com/gin-gonic/gin"
)

type Server struct {
	db     *sql.DB
	config *config.Config
	router *gin.Engine
}

func NewServer(db *sql.DB, cfg *config.Config) *Server {
	server := &Server{
		db:     db,
		config: cfg,
		router: gin.Default(),
	}

	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	// Health check endpoint
	s.router.GET("/health", s.healthCheck)

	// API v1 routes
	v1 := s.router.Group("/api/v1")
	{
		// Public routes
		v1.POST("/auth/login", s.login)
		v1.POST("/server/auth", s.serverAuth)
		
		// Protected routes
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(s.config.JWT.Secret))
		{
			// Player routes
			protected.GET("/player/profile", s.getPlayerProfile)
			protected.PUT("/player/profile", s.updatePlayerProfile)
			protected.GET("/player/stats", s.getPlayerStats)
			protected.PUT("/player/stats", s.updatePlayerStats)
			protected.GET("/player/data/:key", s.getPlayerData)
			protected.PUT("/player/data/:key", s.setPlayerData)
			
			// Pokédex routes
			protected.GET("/pokedex/summary", s.getPokedexSummary)
			protected.POST("/pokedex/region", s.getRegionalPokedex)
			protected.PUT("/pokedex/update", s.updatePokedex)
			protected.PUT("/pokedex/catch", s.updatePokedexSimple) // Simplified national dex endpoint
			protected.GET("/pokedex/leaderboard", s.getPokedexLeaderboard)
		}
		
		// Server proxy routes (for Minecraft server communication)
		server := v1.Group("/server")
		server.Use(middleware.ServerAuthMiddleware(s.config.JWT.Secret))
		{
			// Player management
			server.POST("/player/get", s.serverGetPlayer)
			server.POST("/player/create", s.serverCreateOrUpdatePlayer)
			server.POST("/player/stats/get", s.serverGetPlayerStats)
			server.POST("/player/stats/update", s.serverUpdatePlayerStats)
			server.POST("/player/data/get", s.serverGetPlayerData)
			server.POST("/player/data/set", s.serverSetPlayerData)
			
			// Pokédex management
			server.POST("/pokedex/summary", s.serverGetPokedexSummary)
			server.POST("/pokedex/region", s.serverGetRegionalPokedex)
			server.POST("/pokedex/update", s.serverUpdatePokedex)
			server.GET("/pokedex/leaderboard", s.getPokedexLeaderboard)
		}
		
		// Web dashboard routes (public - for web frontend)
		web := v1.Group("/web")
		{
			web.GET("/leaderboards", s.getWebLeaderboards)
			web.GET("/player/:username/stats", s.getWebPlayerStats)
			web.GET("/server/analytics", s.getWebServerAnalytics)
			web.GET("/pokemon/:dex/popularity", s.getWebPokemonPopularity)
		}
	}
}

func (s *Server) Run(addr string) error {
	return s.router.Run(addr)
}