package auth

import (
	"errors"
	"strconv"
	"time"

	"servixa/internal/model"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID uint           `json:"user_id"`
	Role   model.UserRole `json:"role"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(
	userID uint,
	role model.UserRole,
	secret string,
	duration time.Duration,
) (string, error) {
	now := time.Now()

	claims := JWTClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatUint(uint64(userID), 10),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(duration)),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		claims,
	)

	return token.SignedString([]byte(secret))
}

func ParseAccessToken(
	tokenString string,
	secret string,
) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&JWTClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, errors.New("unexpected signing method")
			}

			return []byte(secret), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid access token")
	}

	return claims, nil
}
