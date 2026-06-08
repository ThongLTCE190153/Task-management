package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte
var jwtExpireHours int

func SetJWTConfig(secret string, expireHours int) {
	jwtSecret = []byte(secret)
	jwtExpireHours = expireHours
}

func GenerateJWT(userID int, role string) (string, error) {
	if len(jwtSecret) == 0 {
		return "", errors.New("jwt secret is not configured")
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(time.Duration(jwtExpireHours) * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtSecret)
}

func VerifyJWT(tokenString string) (*jwt.Token, error) {
	if len(jwtSecret) == 0 {
		return nil, errors.New("jwt secret is not configured")
	}

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		return jwtSecret, nil
	})
}


func GenerateRefreshToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"type":    "refresh", // ← phân biệt với access token
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(), // hết hạn sau 7 ngày
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(jwtSecret))
}

func VerifyRefreshToken(tokenString string) (*jwt.Token, error) {
	return VerifyJWT(tokenString) // dùng lại hàm verify cũ
}
