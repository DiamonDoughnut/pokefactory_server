package api

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (s *Server) generateServerToken(serverID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"server_id": serverID,
		"type":      "server",
		"exp":       time.Now().Add(time.Hour * 24 * 30).Unix(), // 30 days for server tokens
	})

	return token.SignedString([]byte(s.config.JWT.Secret))
}