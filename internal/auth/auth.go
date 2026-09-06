package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("secret-key")

func GenerateToken(userID int) (string, error) {
	claims := jwt.MapClaims{"user_id": userID, "exp": time.Now().Add(time.Hour * 24).Unix()}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
func ParseToken(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || !token.Valid {
		return 0, errors.New("Invalid token")
	}
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["user_id"].(float64)
	return int(userID), nil
}
