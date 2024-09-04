package jwt

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/tgkzz/auth/internal/domain/models"
	"time"
)

const appSecret = "kamal_secret"

func NewToken(user models.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["username"] = user.Username
	claims["exp"] = time.Now().Add(24 * time.Hour).Unix()

	tokenString, err := token.SignedString([]byte(appSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
