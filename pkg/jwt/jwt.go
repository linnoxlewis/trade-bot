package jwt

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"time"
)

var (
	errNullClaims   = errors.New("token or claims is null")
	errParseToken   = errors.New("can't decode token")
	errInvalidToken = errors.New("invalid token")
)

type JwtClaims struct {
	UserID uuid.UUID `json:"uid"`
	jwt.StandardClaims
}

func NewJwtClaims(userID uuid.UUID, ttl time.Duration) *JwtClaims {
	now := time.Now().UTC()
	return &JwtClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  now.Unix(),
			ExpiresAt: now.Add(ttl).Unix(),
		},
	}
}

func NewTokenString(userID uuid.UUID, ttl time.Duration, secret string) (string, error) {
	claims := NewJwtClaims(userID, ttl)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseToken(inputToken, secret string) (*JwtClaims, error) {
	token, err := jwt.ParseWithClaims(inputToken, &JwtClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errParseToken
		}

		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	if token == nil || token.Claims == nil {
		return nil, errNullClaims
	}

	claims, ok := token.Claims.(*JwtClaims)
	if !ok || !token.Valid {
		return nil, errInvalidToken
	}

	return claims, nil
}
