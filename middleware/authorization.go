package middleware

import (
	"golang-final-project/config"

	"github.com/gin-gonic/gin"
)

func isAccessTokenAssigned(token string) bool {
	var exists bool

	query := `
	SELECT EXISTS (
		SELECT 1 FROM users WHERE token = $1
	)`

	err := config.Db.QueryRow(
		query,
		token,
	).Scan(&exists)

	if err != nil {
		panic(err)
	} else {
		return exists
	}
}

func ValidateAccessToken(ctx *gin.Context) (id string, role string, error string) {
	authHeader := ctx.GetHeader("Authorization")

	if authHeader == "" {
		return "", "", "access token is required"
	} else {
		tokenString := authHeader[len("Bearer "):]

		exists := isAccessTokenAssigned(tokenString)

		if !exists {
			return "", "", "access token is not assigned to any user"
		} else {
			userId, role, err := ValidateJWT(tokenString)

			if err != nil {
				return "", "", "invalid or expired access token"
			} else {
				return userId, role, ""
			}
		}
	}
}
