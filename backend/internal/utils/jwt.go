package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gofrs/uuid"
)

type Claims struct {
	UserID      uuid.UUID     `json:"user_id"`
	Username    string        `json:"username"`
	Roles       []string      `json:"roles"`
	IsAdmin     bool          `json:"is_admin"`
	Permissions []Permission  `json:"permissions"`
	jwt.RegisteredClaims
}

type Permission struct {
	Resource string   `json:"resource"`
	Actions  []string `json:"actions"`
}

var jwtSecret = []byte("your-secret-key-change-this-in-production")

func GenerateJWT(userID uuid.UUID, username string, roles []string, isAdmin bool, permissions []Permission) (string, error) {
	claims := Claims{
		UserID:      userID,
		Username:    username,
		Roles:       roles,
		IsAdmin:     isAdmin,
		Permissions: permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "task-manager",
			Subject:   userID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func ValidateJWT(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func HasPermission(permissions []Permission, resource, action string) bool {
	for _, perm := range permissions {
		if perm.Resource == resource {
			for _, act := range perm.Actions {
				if act == action {
					return true
				}
			}
		}
	}
	return false
}
