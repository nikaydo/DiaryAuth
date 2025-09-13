package jwt

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/nikaydo/DiaryAuth/internal/config"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

var (
	ErrTokenExpired = errors.New("token is expired")
)

type JwtTokens struct {
	AccessToken  string
	RefreshToken string
}

func CreateTokens(uuid uuid.UUID, role string, env config.Env) (JwtTokens, error) {
	var err error
	var tokens JwtTokens
	tokens.AccessToken, err = CreateToken(uuid, role, env.JWTTTL, env.JWTSecret)
	if err != nil {
		return tokens, fmt.Errorf("error creating JWT token: %w", err)
	}
	tokens.RefreshToken, err = CreateToken(uuid, role, env.RefreshTTL, env.RefreshSecret)
	if err != nil {
		return tokens, fmt.Errorf("error creating refresh token: %w", err)
	}
	return tokens, nil
}

func CreateToken(uuid uuid.UUID, role, tokenTTL, secret string) (string, error) {
	exTime, err := strconv.Atoi(tokenTTL)
	if err != nil {
		return "", fmt.Errorf("failed to parse TTL from environment: %w", err)
	}
	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub":  uuid,
			"iss":  "diary",
			"role": role,
			"exp":  time.Now().Add(time.Duration(exTime * int(time.Minute))).Unix(),
			"iat":  time.Now().Unix(),
		}).SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("error signing token: %w", err)
	}
	return tokenString, nil
}

func ValidateToken(jwtToken, secret string) (uuid.UUID, error) {
	token, err := jwt.Parse(jwtToken, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("no valid signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		if err.Error() == "Token is expired" {
			uuid, err := setClaims(token)
			if err != nil {
				return uuid, fmt.Errorf("failed to extract claims from expired token: %w", err)
			}
			return uuid, ErrTokenExpired
		}
		return uuid.Nil, fmt.Errorf("failed to parse token: %w", err)
	}
	return setClaims(token)
}

func setClaims(token *jwt.Token) (uuid.UUID, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, fmt.Errorf("invalid token")
	}
	id, ok := claims["sub"].(string)
	if !ok {
		return uuid.Nil, fmt.Errorf("cant parse id from jwt token")
	}
	return uuid.MustParse(id), nil
}
