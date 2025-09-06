package models

import (
	"time"
)

type Player struct {
	ID           int       `json:"id" db:"id"`
	UUID         string    `json:"uuid" db:"uuid"`
	Username     string    `json:"username" db:"username"`
	LastLogin    time.Time `json:"last_login" db:"last_login"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type PlayerStats struct {
	ID           int     `json:"id" db:"id"`
	PlayerID     int     `json:"player_id" db:"player_id"`
	Level        int     `json:"level" db:"level"`
	Experience   int     `json:"experience" db:"experience"`
	Currency     int     `json:"currency" db:"currency"`
	PlayTime     int     `json:"play_time" db:"play_time"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

type PlayerData struct {
	ID           int       `json:"id" db:"id"`
	PlayerID     int       `json:"player_id" db:"player_id"`
	DataKey      string    `json:"data_key" db:"data_key"`
	DataValue    string    `json:"data_value" db:"data_value"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}