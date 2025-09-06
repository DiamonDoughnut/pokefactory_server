package models

import (
	"time"
)

type PokedexSummary struct {
	ID                         int       `json:"id" db:"id"`
	PlayerID                   int       `json:"player_id" db:"player_id"`
	TotalCaught                int       `json:"total_caught" db:"total_caught"`
	TotalSeen                  int       `json:"total_seen" db:"total_seen"`
	RegionsCompleted           int       `json:"regions_completed" db:"regions_completed"`
	NationalCompletionPercent  float64   `json:"national_completion_percentage" db:"national_completion_percentage"`
	LastUpdated                time.Time `json:"last_updated" db:"last_updated"`
	CreatedAt                  time.Time `json:"created_at" db:"created_at"`
}

type RegionalPokedex struct {
	ID                   int       `json:"id" db:"id"`
	PlayerID             int       `json:"player_id" db:"player_id"`
	CaughtFlags          []byte    `json:"caught_flags" db:"caught_flags"`
	SeenFlags            []byte    `json:"seen_flags" db:"seen_flags"`
	CompletionDate       *time.Time `json:"completion_date" db:"completion_date"`
	CompletionPercentage float64   `json:"completion_percentage" db:"completion_percentage"`
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" db:"updated_at"`
}

type PokedexUpdateRequest struct {
	Region       string `json:"region,omitempty"`        // Optional - will be auto-determined if not provided
	PokemonID    int    `json:"pokemon_id" binding:"required"` // Can be national or regional ID
	NationalID   int    `json:"national_id,omitempty"`   // Optional - use this for national dex numbers
	Action       string `json:"action" binding:"required"` // "catch" or "see"
}

type LeaderboardEntry struct {
	PlayerID                  int     `json:"player_id" db:"player_id"`
	Username                  string  `json:"username" db:"username"`
	NationalCompletionPercent float64 `json:"national_completion_percentage" db:"national_completion_percentage"`
	TotalCaught               int     `json:"total_caught" db:"total_caught"`
}