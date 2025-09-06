package api

import (
	"fmt"

	"pokefactory_server/internal/models"
)

var regionTables = map[string]string{
	"kanto":  "player_pokedex_kanto",
	"johto":  "player_pokedex_johto",
	"hoenn":  "player_pokedex_hoenn",
	"sinnoh": "player_pokedex_sinnoh",
	"unova":  "player_pokedex_unova",
	"kalos":  "player_pokedex_kalos",
	"alola":  "player_pokedex_alola",
	"galar":  "player_pokedex_galar",
	"hisui":  "player_pokedex_hisui",
	"paldea": "player_pokedex_paldea",
}

var regionSizes = map[string]int{
	"kanto":  151,
	"johto":  100,
	"hoenn":  135,
	"sinnoh": 107,
	"unova":  156,
	"kalos":  72,
	"alola":  88,
	"galar":  89,
	"hisui":  7,
	"paldea": 120,
}

// National dex number ranges for each region
var nationalDexRanges = map[string][2]int{
	"kanto":  {1, 151},
	"johto":  {152, 251},
	"hoenn":  {252, 386},
	"sinnoh": {387, 493},
	"unova":  {494, 649},
	"kalos":  {650, 721},
	"alola":  {722, 809},
	"galar":  {810, 898},
	"hisui":  {899, 905},
	"paldea": {906, 1008},
}

func (s *Server) getOrCreatePokedexSummary(playerID int) (*models.PokedexSummary, error) {
	summary, err := s.getPokedexSummaryByID(playerID)
	if err == nil {
		return summary, nil
	}

	query := `
		INSERT INTO player_pokedex_summary (player_id, created_at)
		VALUES ($1, NOW())
		RETURNING id, player_id, total_caught, total_seen, regions_completed, national_completion_percentage, last_updated, created_at`

	summary = &models.PokedexSummary{}
	err = s.db.QueryRow(query, playerID).Scan(
		&summary.ID, &summary.PlayerID, &summary.TotalCaught, &summary.TotalSeen,
		&summary.RegionsCompleted, &summary.NationalCompletionPercent,
		&summary.LastUpdated, &summary.CreatedAt,
	)

	return summary, err
}

func (s *Server) getPokedexSummaryByID(playerID int) (*models.PokedexSummary, error) {
	query := `SELECT id, player_id, total_caught, total_seen, regions_completed, national_completion_percentage, last_updated, created_at FROM player_pokedex_summary WHERE player_id = $1`
	
	summary := &models.PokedexSummary{}
	err := s.db.QueryRow(query, playerID).Scan(
		&summary.ID, &summary.PlayerID, &summary.TotalCaught, &summary.TotalSeen,
		&summary.RegionsCompleted, &summary.NationalCompletionPercent,
		&summary.LastUpdated, &summary.CreatedAt,
	)
	
	return summary, err
}

func (s *Server) getOrCreateRegionalPokedex(playerID int, region string) (*models.RegionalPokedex, error) {
	pokedex, err := s.getRegionalPokedexByID(playerID, region)
	if err == nil {
		return pokedex, nil
	}

	tableName, exists := regionTables[region]
	if !exists {
		return nil, fmt.Errorf("invalid region: %s", region)
	}

	query := fmt.Sprintf(`
		INSERT INTO %s (player_id, created_at, updated_at)
		VALUES ($1, NOW(), NOW())
		RETURNING id, player_id, caught_flags, seen_flags, completion_date, completion_percentage, created_at, updated_at`, tableName)

	pokedex = &models.RegionalPokedex{}
	err = s.db.QueryRow(query, playerID).Scan(
		&pokedex.ID, &pokedex.PlayerID, &pokedex.CaughtFlags, &pokedex.SeenFlags,
		&pokedex.CompletionDate, &pokedex.CompletionPercentage,
		&pokedex.CreatedAt, &pokedex.UpdatedAt,
	)

	return pokedex, err
}

func (s *Server) getRegionalPokedexByID(playerID int, region string) (*models.RegionalPokedex, error) {
	tableName, exists := regionTables[region]
	if !exists {
		return nil, fmt.Errorf("invalid region: %s", region)
	}

	query := fmt.Sprintf(`SELECT id, player_id, caught_flags, seen_flags, completion_date, completion_percentage, created_at, updated_at FROM %s WHERE player_id = $1`, tableName)
	
	pokedex := &models.RegionalPokedex{}
	err := s.db.QueryRow(query, playerID).Scan(
		&pokedex.ID, &pokedex.PlayerID, &pokedex.CaughtFlags, &pokedex.SeenFlags,
		&pokedex.CompletionDate, &pokedex.CompletionPercentage,
		&pokedex.CreatedAt, &pokedex.UpdatedAt,
	)
	
	return pokedex, err
}

func (s *Server) updatePokedexEntry(playerID int, req models.PokedexUpdateRequest) error {
	var region string
	var regionalID int
	var err error

	// Determine region and regional ID
	if req.NationalID > 0 {
		// Use national dex number to find region
		region, regionalID, err = getRegionFromNationalDex(req.NationalID)
		if err != nil {
			return err
		}
	} else if req.Region != "" {
		// Use provided region and pokemon_id as regional ID
		region = req.Region
		regionalID = req.PokemonID
	} else {
		return fmt.Errorf("either region+pokemon_id or national_id must be provided")
	}

	// Get or create regional pokedex
	pokedex, err := s.getOrCreateRegionalPokedex(playerID, region)
	if err != nil {
		return err
	}

	// Update bitfield based on action
	var updatedFlags []byte
	if req.Action == "catch" {
		updatedFlags = setBit(pokedex.CaughtFlags, regionalID-1)
		if err := s.updateRegionalFlags(playerID, region, "caught_flags", updatedFlags); err != nil {
			return err
		}
	} else if req.Action == "see" {
		updatedFlags = setBit(pokedex.SeenFlags, regionalID-1)
		if err := s.updateRegionalFlags(playerID, region, "seen_flags", updatedFlags); err != nil {
			return err
		}
	}

	// Update completion percentage and summary
	return s.updatePokedexCompletion(playerID, region)
}

func (s *Server) updateRegionalFlags(playerID int, region, flagType string, flags []byte) error {
	tableName := regionTables[region]
	query := fmt.Sprintf(`UPDATE %s SET %s = $1, updated_at = NOW() WHERE player_id = $2`, tableName, flagType)
	_, err := s.db.Exec(query, flags, playerID)
	return err
}

func (s *Server) updatePokedexCompletion(playerID int, region string) error {
	// Calculate regional completion
	pokedex, err := s.getRegionalPokedexByID(playerID, region)
	if err != nil {
		return err
	}

	regionSize := regionSizes[region]
	caughtCount := countBits(pokedex.CaughtFlags)
	completionPercent := float64(caughtCount) / float64(regionSize) * 100

	// Update regional completion
	tableName := regionTables[region]
	query := fmt.Sprintf(`UPDATE %s SET completion_percentage = $1, updated_at = NOW() WHERE player_id = $2`, tableName)
	if _, err := s.db.Exec(query, completionPercent, playerID); err != nil {
		return err
	}

	// Update summary
	return s.updatePokedexSummaryStats(playerID)
}

func (s *Server) updatePokedexSummaryStats(playerID int) error {
	// Calculate totals across all regions
	totalCaught := 0
	totalSeen := 0
	regionsCompleted := 0

	for region := range regionTables {
		pokedex, err := s.getRegionalPokedexByID(playerID, region)
		if err != nil {
			continue
		}
		
		caughtCount := countBits(pokedex.CaughtFlags)
		seenCount := countBits(pokedex.SeenFlags)
		
		totalCaught += caughtCount
		totalSeen += seenCount
		
		if pokedex.CompletionPercentage >= 100.0 {
			regionsCompleted++
		}
	}

	nationalPercent := float64(totalCaught) / 1018.0 * 100 // Total Pokémon across all regions

	query := `
		UPDATE player_pokedex_summary 
		SET total_caught = $1, total_seen = $2, regions_completed = $3, 
		    national_completion_percentage = $4, last_updated = NOW()
		WHERE player_id = $5`
	
	_, err := s.db.Exec(query, totalCaught, totalSeen, regionsCompleted, nationalPercent, playerID)
	return err
}

func (s *Server) getPokedexLeaderboardData() ([]models.LeaderboardEntry, error) {
	query := `
		SELECT ps.player_id, p.username, ps.national_completion_percentage, ps.total_caught
		FROM player_pokedex_summary ps
		JOIN players p ON ps.player_id = p.id
		ORDER BY ps.national_completion_percentage DESC, ps.total_caught DESC
		LIMIT 50`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var leaderboard []models.LeaderboardEntry
	for rows.Next() {
		var entry models.LeaderboardEntry
		if err := rows.Scan(&entry.PlayerID, &entry.Username, &entry.NationalCompletionPercent, &entry.TotalCaught); err != nil {
			continue
		}
		leaderboard = append(leaderboard, entry)
	}

	return leaderboard, nil
}

func setBit(data []byte, position int) []byte {
	if len(data) == 0 {
		data = make([]byte, (position/8)+1)
	}
	
	byteIndex := position / 8
	bitIndex := position % 8
	
	if byteIndex >= len(data) {
		newData := make([]byte, byteIndex+1)
		copy(newData, data)
		data = newData
	}
	
	data[byteIndex] |= (1 << bitIndex)
	return data
}

func countBits(data []byte) int {
	count := 0
	for _, b := range data {
		for b != 0 {
			count += int(b & 1)
			b >>= 1
		}
	}
	return count
}

func getRegionFromNationalDex(nationalID int) (string, int, error) {
	for region, rangeData := range nationalDexRanges {
		if nationalID >= rangeData[0] && nationalID <= rangeData[1] {
			regionalID := nationalID - rangeData[0] + 1
			return region, regionalID, nil
		}
	}
	return "", 0, fmt.Errorf("national dex number %d not found in any region", nationalID)
}