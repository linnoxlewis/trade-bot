package pkg

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"time"
)

var (
	errNullClaims   = errors.New("err token or claims is null")
	errParseToken   = errors.New("err can't decode token")
	errInvalidToken = errors.New("err no valid token")
)

type Jwt struct {
	UserID uuid.UUID `json:"uid"`
	jwt.StandardClaims
}

func NewJwt(userID uuid.UUID, ttl time.Duration) (*Jwt, error) {
	now := time.Now()
	return &Jwt{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  now.UTC().Unix(),
			ExpiresAt: now.Add(ttl).UTC().Unix(),
		}}, nil
}

func NewTokenString(userID uuid.UUID, ttl time.Duration, secret string) (string, error) {
	auth, err := NewJwt(userID, ttl)
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, auth)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseToken(token, secret string) (*Jwt, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &Jwt{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errParseToken
		}

		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	if jwtToken == nil || jwtToken.Claims == nil {
		return nil, errNullClaims
	}

	authJwt, ok := jwtToken.Claims.(*Jwt)
	if !ok || !jwtToken.Valid {
		return nil, errInvalidToken
	}

	return authJwt, nil
}
