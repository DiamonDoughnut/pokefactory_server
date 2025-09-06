package models

import (
	"time"
)

type WebLeaderboard struct {
	PlayerID                  int     `json:"player_id" db:"player_id"`
	Username                  string  `json:"username" db:"username"`
	Level                     int     `json:"level" db:"level"`
	TotalCaught               int     `json:"total_caught" db:"total_caught"`
	NationalCompletionPercent float64 `json:"national_completion_percentage" db:"national_completion_percentage"`
	LastLogin                 time.Time `json:"last_login" db:"last_login"`
}

type WebPlayerStats struct {
	Player   Player          `json:"player"`
	Stats    PlayerStats     `json:"stats"`
	Pokedex  PokedexSummary  `json:"pokedex"`
}

type WebServerAnalytics struct {
	TotalPlayers        int     `json:"total_players"`
	ActivePlayers       int     `json:"active_players_24h"`
	TotalPokemonCaught  int     `json:"total_pokemon_caught"`
	AverageLevel        float64 `json:"average_level"`
	TopRegion           string  `json:"most_popular_region"`
	ServerUptime        string  `json:"server_uptime"`
}

type WebPokemonPopularity struct {
	NationalID    int     `json:"national_id"`
	Name          string  `json:"name,omitempty"`
	Region        string  `json:"region"`
	CatchCount    int     `json:"catch_count"`
	SeenCount     int     `json:"seen_count"`
	PopularityRank int    `json:"popularity_rank"`
	CatchRate     float64 `json:"catch_rate_percentage"`
}