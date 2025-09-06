package models

// Server-to-API request models for Minecraft server proxy operations
type ServerPlayerRequest struct {
	PlayerUUID string `json:"player_uuid" binding:"required"`
}

type ServerPlayerUpdateRequest struct {
	PlayerUUID string `json:"player_uuid" binding:"required"`
	Username   string `json:"username"`
}

type ServerPlayerStatsRequest struct {
	PlayerUUID string      `json:"player_uuid" binding:"required"`
	Stats      PlayerStats `json:"stats,omitempty"`
}

type ServerPlayerDataRequest struct {
	PlayerUUID string `json:"player_uuid" binding:"required"`
	DataKey    string `json:"data_key" binding:"required"`
	DataValue  string `json:"data_value,omitempty"`
}

type ServerPokedexRequest struct {
	PlayerUUID string `json:"player_uuid" binding:"required"`
	Region     string `json:"region,omitempty"`
}

type ServerPokedexUpdateRequest struct {
	PlayerUUID string `json:"player_uuid" binding:"required"`
	NationalID int    `json:"national_id" binding:"required"`
	Action     string `json:"action" binding:"required"` // "catch" or "see"
}

type ServerAuthRequest struct {
	ServerID string `json:"server_id" binding:"required"`
	ServerKey string `json:"server_key" binding:"required"`
}