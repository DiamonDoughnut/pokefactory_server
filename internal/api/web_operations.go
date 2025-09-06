package api

import (
	"fmt"

	"pokefactory_server/internal/models"
)

func (s *Server) getWebLeaderboardData() ([]models.WebLeaderboard, error) {
	query := `
		SELECT p.id, p.username, ps.level, pds.total_caught, 
		       pds.national_completion_percentage, p.last_login
		FROM players p
		JOIN player_stats ps ON p.id = ps.player_id
		LEFT JOIN player_pokedex_summary pds ON p.id = pds.player_id
		ORDER BY pds.national_completion_percentage DESC, ps.level DESC, pds.total_caught DESC
		LIMIT 100`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var leaderboards []models.WebLeaderboard
	for rows.Next() {
		var lb models.WebLeaderboard
		var totalCaught *int
		var completionPercentFloat *float64

		err := rows.Scan(&lb.PlayerID, &lb.Username, &lb.Level, 
			&totalCaught, &completionPercentFloat, &lb.LastLogin)
		if err != nil {
			continue
		}

		if totalCaught != nil {
			lb.TotalCaught = *totalCaught
		}
		if completionPercentFloat != nil {
			lb.NationalCompletionPercent = *completionPercentFloat
		}

		leaderboards = append(leaderboards, lb)
	}

	return leaderboards, nil
}

func (s *Server) getServerAnalyticsData() (*models.WebServerAnalytics, error) {
	analytics := &models.WebServerAnalytics{}

	// Total players
	err := s.db.QueryRow("SELECT COUNT(*) FROM players").Scan(&analytics.TotalPlayers)
	if err != nil {
		return nil, err
	}

	// Active players (last 24 hours)
	err = s.db.QueryRow(`
		SELECT COUNT(*) FROM players 
		WHERE last_login > NOW() - INTERVAL '24 hours'`).Scan(&analytics.ActivePlayers)
	if err != nil {
		analytics.ActivePlayers = 0
	}

	// Total Pokemon caught across all players
	err = s.db.QueryRow(`
		SELECT COALESCE(SUM(total_caught), 0) FROM player_pokedex_summary`).Scan(&analytics.TotalPokemonCaught)
	if err != nil {
		analytics.TotalPokemonCaught = 0
	}

	// Average level
	err = s.db.QueryRow(`
		SELECT COALESCE(AVG(level), 0) FROM player_stats`).Scan(&analytics.AverageLevel)
	if err != nil {
		analytics.AverageLevel = 0
	}

	// Most popular region (highest completion rates)
	analytics.TopRegion = s.getMostPopularRegion()

	// Server uptime (placeholder - would need actual server start time)
	analytics.ServerUptime = "Available via /health endpoint"

	return analytics, nil
}

func (s *Server) getMostPopularRegion() string {
	regions := []string{"kanto", "johto", "hoenn", "sinnoh", "unova", "alola", "galar", "hisui", "paldea"}
	maxCompletion := 0.0
	topRegion := "kanto"

	for _, region := range regions {
		tableName := regionTables[region]
		var avgCompletion float64
		
		query := fmt.Sprintf(`
			SELECT COALESCE(AVG(completion_percentage), 0) 
			FROM %s WHERE completion_percentage > 0`, tableName)
		
		err := s.db.QueryRow(query).Scan(&avgCompletion)
		if err == nil && avgCompletion > maxCompletion {
			maxCompletion = avgCompletion
			topRegion = region
		}
	}

	return topRegion
}

func (s *Server) getPokemonPopularityData(nationalID int) (*models.WebPokemonPopularity, error) {
	region, regionalID, err := getRegionFromNationalDex(nationalID)
	if err != nil {
		return nil, err
	}

	popularity := &models.WebPokemonPopularity{
		NationalID: nationalID,
		Region:     region,
	}

	// Count how many players have caught this Pokemon
	tableName := regionTables[region]
	bitPosition := regionalID - 1

	// This is a simplified query - in production you'd use proper bitwise operations
	query := fmt.Sprintf(`
		SELECT COUNT(*) as catch_count,
		       (SELECT COUNT(*) FROM %s WHERE LENGTH(caught_flags) > %d) as seen_count,
		       (SELECT COUNT(*) FROM %s) as total_players
		FROM %s 
		WHERE LENGTH(caught_flags) > %d`, 
		tableName, bitPosition/8, tableName, tableName, bitPosition/8)

	var totalPlayers int
	err = s.db.QueryRow(query).Scan(&popularity.CatchCount, &popularity.SeenCount, &totalPlayers)
	if err != nil {
		return nil, err
	}

	// Calculate catch rate percentage
	if totalPlayers > 0 {
		popularity.CatchRate = float64(popularity.CatchCount) / float64(totalPlayers) * 100
	}

	// Get popularity rank (simplified - would need more complex ranking logic)
	popularity.PopularityRank = s.calculatePopularityRank(nationalID, popularity.CatchCount)

	return popularity, nil
}

func (s *Server) calculatePopularityRank(nationalID, catchCount int) int {
	// Simplified ranking - in production would rank against all Pokemon
	var rank int
	query := `
		SELECT COUNT(*) + 1 FROM (
			SELECT 1 FROM player_pokedex_summary 
			WHERE total_caught > $1
		) as higher_counts`
	
	err := s.db.QueryRow(query, catchCount).Scan(&rank)
	if err != nil {
		return 999 // Default rank if calculation fails
	}
	
	return rank
}