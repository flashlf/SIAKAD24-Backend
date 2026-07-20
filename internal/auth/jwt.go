package auth

import (
	"errors"
	"os"
	"sync"

	"github.com/golang-jwt/jwt/v5"
)

type Role string

const (
	RoleAdministrator Role = "Administrator"
	RoleStaff         Role = "Staff"
)

func (r Role) valid() bool {
	return r == RoleAdministrator || r == RoleStaff
}

type Claims struct {
	Subject string `json:"sub"`
	Role    Role   `json:"role"`
	jwt.RegisteredClaims
}

var (
	jwtSecretOnce sync.Once
	jwtSecret     []byte
)

func secret() []byte {
	jwtSecretOnce.Do(func() {
		jwtSecret = []byte(os.Getenv("JWT_SECRET"))
	})
	return jwtSecret
}

var ErrInvalidToken = errors.New("invalid token")

// ParseToken validates the HS256 signature and expiry of tokenString and
// extracts its custom claims. A token whose "role" claim is not Administrator
// or Staff is rejected as invalid, same as a bad signature or expired token.
func ParseToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return secret(), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, ErrInvalidToken
	}
	if !claims.Role.valid() {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
