package authrepository

import (
	"fmt"
	"github.com/Algoru/frontera/configuration"
	"github.com/badoux/checkmail"
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

func (a *Auth) HasRequiredFields() []string {
	errors := make([]string, 0)

	if err := checkmail.ValidateFormat(a.Email); err != nil {
		str := fmt.Sprintf("invalid email: %s", err.Error())
		errors = append(errors, str)
	}

	minPasswordLength := int(configuration.GetConfiguration().Security.MinPasswordLength)
	if len(a.Password) < minPasswordLength {
		str := fmt.Sprintf("the password must contain at least %d characters", minPasswordLength)
		errors = append(errors, str)
	}

	if len(errors) == 0 {
		return nil
	}

	return errors
}