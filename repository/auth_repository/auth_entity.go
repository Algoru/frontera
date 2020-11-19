package authrepository

import (
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/microcosm-cc/bluemonday"
)

// Auth
type Auth struct {
	Email    string
	Password string
}

// AuthClaims
type AuthClaims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

// Sanitize
func (a *Auth) Sanitize() {
	a.Email = strings.ToLower(a.Email)
	a.Email = bluemonday.StrictPolicy().Sanitize(a.Email)
}
