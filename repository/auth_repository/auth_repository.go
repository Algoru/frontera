package authrepository

import (
	"github.com/Algoru/frontera/domain/entity"
)

// AuthRepository
type AuthRepository interface {
	GetCredentialByToken(string) (*entity.Credential, error)
	AddUserSession(*entity.Credential) error
	RemoveUserSessions(string) error
	RemoveSingleSession(string, string) error
}
