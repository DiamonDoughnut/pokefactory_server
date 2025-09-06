package api

import (
	"database/sql"
	"time"

	"pokefactory_server/internal/models"
)

func (s *Server) getOrCreatePlayer(uuid, username string) (*models.Player, error) {
	// Try to get existing player
	player, err := s.getPlayerByUUID(uuid)
	if err == nil {
		// Update last login and username
		s.updatePlayerLogin(player.ID, username)
		return player, nil
	}

	// Create new player
	query := `
		INSERT INTO players (uuid, username, last_login, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW(), NOW())
		RETURNING id, uuid, username, last_login, created_at, updated_at`

	player = &models.Player{}
	err = s.db.QueryRow(query, uuid, username).Scan(
		&player.ID, &player.UUID, &player.Username,
		&player.LastLogin, &player.CreatedAt, &player.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Create initial stats and pokedex
	s.createPlayerStats(player.ID)
	s.getOrCreatePokedexSummary(player.ID)

	return player, nil
}

func (s *Server) getPlayerByUUID(uuid string) (*models.Player, error) {
	query := `SELECT id, uuid, username, last_login, created_at, updated_at FROM players WHERE uuid = $1`
	
	player := &models.Player{}
	err := s.db.QueryRow(query, uuid).Scan(
		&player.ID, &player.UUID, &player.Username,
		&player.LastLogin, &player.CreatedAt, &player.UpdatedAt,
	)
	
	return player, err
}

func (s *Server) getPlayerByUsername(username string) (*models.Player, error) {
	query := `SELECT id, uuid, username, last_login, created_at, updated_at FROM players WHERE username = $1`
	
	player := &models.Player{}
	err := s.db.QueryRow(query, username).Scan(
		&player.ID, &player.UUID, &player.Username,
		&player.LastLogin, &player.CreatedAt, &player.UpdatedAt,
	)
	
	return player, err
}

func (s *Server) updatePlayerLogin(playerID int, username string) error {
	query := `UPDATE players SET username = $1, last_login = NOW(), updated_at = NOW() WHERE id = $2`
	_, err := s.db.Exec(query, username, playerID)
	return err
}

func (s *Server) updatePlayer(playerID int, username string) error {
	query := `UPDATE players SET username = $1, updated_at = NOW() WHERE id = $2`
	_, err := s.db.Exec(query, username, playerID)
	return err
}

func (s *Server) createPlayerStats(playerID int) error {
	query := `
		INSERT INTO player_stats (player_id, level, experience, currency, play_time, created_at, updated_at)
		VALUES ($1, 1, 0, 0, 0, NOW(), NOW())`
	_, err := s.db.Exec(query, playerID)
	return err
}

func (s *Server) getPlayerStatsByID(playerID int) (*models.PlayerStats, error) {
	query := `SELECT id, player_id, level, experience, currency, play_time, created_at, updated_at FROM player_stats WHERE player_id = $1`
	
	stats := &models.PlayerStats{}
	err := s.db.QueryRow(query, playerID).Scan(
		&stats.ID, &stats.PlayerID, &stats.Level, &stats.Experience,
		&stats.Currency, &stats.PlayTime, &stats.CreatedAt, &stats.UpdatedAt,
	)
	
	return stats, err
}

func (s *Server) updatePlayerStatsByID(stats models.PlayerStats) error {
	query := `
		UPDATE player_stats 
		SET level = $1, experience = $2, currency = $3, play_time = $4, updated_at = NOW()
		WHERE player_id = $5`
	_, err := s.db.Exec(query, stats.Level, stats.Experience, stats.Currency, stats.PlayTime, stats.PlayerID)
	return err
}

func (s *Server) getPlayerDataByKey(playerID int, key string) (*models.PlayerData, error) {
	query := `SELECT id, player_id, data_key, data_value, created_at, updated_at FROM player_data WHERE player_id = $1 AND data_key = $2`
	
	data := &models.PlayerData{}
	err := s.db.QueryRow(query, playerID, key).Scan(
		&data.ID, &data.PlayerID, &data.DataKey, &data.DataValue,
		&data.CreatedAt, &data.UpdatedAt,
	)
	
	return data, err
}

func (s *Server) setPlayerDataByKey(playerID int, key, value string) error {
	query := `
		INSERT INTO player_data (player_id, data_key, data_value, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		ON CONFLICT (player_id, data_key)
		DO UPDATE SET data_value = $3, updated_at = NOW()`
	_, err := s.db.Exec(query, playerID, key, value)
	return err
}