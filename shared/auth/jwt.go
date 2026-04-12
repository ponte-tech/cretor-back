package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UsuarioID string `json:"usuario_id"`
	TenantID  string `json:"tenant_id"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	jwt.RegisteredClaims
}

type JWTProvider struct {
	secret             string
	accessTokenExpiry  time.Duration
	refreshTokenExpiry time.Duration
}

func NewJWTProvider(secret string, accessHours, refreshHours int) *JWTProvider {
	return &JWTProvider{
		secret:             secret,
		accessTokenExpiry:  time.Duration(accessHours) * time.Hour,
		refreshTokenExpiry: time.Duration(refreshHours) * time.Hour,
	}
}

func (p *JWTProvider) GenerateAccessToken(usuarioID, tenantID, email, role string) (string, int64, error) {
	now := time.Now().UTC()
	exp := now.Add(p.accessTokenExpiry)

	claims := &Claims{
		UsuarioID: usuarioID,
		TenantID:  tenantID,
		Email:     email,
		Role:      role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(p.secret))
	if err != nil {
		return "", 0, err
	}
	return signed, int64(p.accessTokenExpiry.Seconds()), nil
}

func (p *JWTProvider) GenerateRefreshToken(usuarioID, tenantID string) (string, error) {
	now := time.Now().UTC()

	claims := &Claims{
		UsuarioID: usuarioID,
		TenantID:  tenantID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(p.refreshTokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(p.secret))
}

func (p *JWTProvider) ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(p.secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
