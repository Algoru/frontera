package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/Algoru/frontera/configuration"

	"github.com/Algoru/frontera/domain/entity"
	authrepository "github.com/Algoru/frontera/repository/auth_repository"
	userrepository "github.com/Algoru/frontera/repository/user_repository"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(*authrepository.Auth) (*entity.Credential, error)
	Logout(string, bool) error
	RemoveUserSessions(string) error
	AddUserSession(*entity.Credential) error
	GetCredentialByToken(string) (*entity.Credential, error)
	RemoveSingleSession(string, string) error
}

type authService struct {
	userRepository userrepository.UserRepository
	authRepository authrepository.AuthRepository
}

func NewAuthService(ur userrepository.UserRepository, ar authrepository.AuthRepository) AuthService {
	return &authService{
		userRepository: ur,
		authRepository: ar,
	}
}

// Login
func (s *authService) Login(auth *authrepository.Auth) (*entity.Credential, error) {
	// 2. Generate JWT
	// 3. Persist to cache
	user, err := s.userRepository.GetUserByEmail(auth.Email)
	if err != nil {
		return nil, err
	}

	userHashedPassword := []byte(user.Password)
	enteredPassword := []byte(auth.Password)
	if err = bcrypt.CompareHashAndPassword(userHashedPassword, enteredPassword); err != nil {
		return nil, err
	}

	token, expiresAt, err := generateToken(user)
	if err != nil {
		return nil, err
	}

	credential := entity.Credential{
		UserID:    user.UserID.String(),
		Token:     token,
		ExpiresAt: expiresAt,
	}

	return &credential, nil
}

func (s *authService) Logout(token string, flush bool) error {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(configuration.GetConfiguration().Security.TokenSigningKey), nil
	})
	if err != nil {
		return err
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return errors.New("unable to parse token claims")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return errors.New("unable to parse user id from token claim")
	}

	if flush {
		err = s.authRepository.RemoveUserSessions(userID)
	} else {
		err = s.authRepository.RemoveSingleSession(userID, token)
	}

	return err
}

func generateToken(user *entity.User) (string, time.Time, error) {
	secConfig := configuration.GetConfiguration().Security

	tokenIssuer := secConfig.TokenIssuer
	if tokenIssuer == "" {
		tokenIssuer = "frontera"
	}

	tokenIssuedAt := time.Now()
	tokenExpiresAt := tokenIssuedAt.Add(time.Hour * time.Duration(secConfig.TokenLifetime))

	authClaims := authrepository.AuthClaims{
		UserID: user.UserID.String(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: tokenExpiresAt.Unix(),
			IssuedAt:  tokenIssuedAt.Unix(),
			Issuer:    tokenIssuer,
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, authClaims)
	token, err := jwtToken.SignedString([]byte(secConfig.TokenSigningKey))
	if err != nil {
		return "", time.Time{}, err
	}

	return token, tokenExpiresAt, nil
}

func (s *authService) RemoveUserSessions(userID string) error {
	return s.authRepository.RemoveUserSessions(userID)
}

func (s *authService) AddUserSession(credential *entity.Credential) error {
	return s.authRepository.AddUserSession(credential)
}

func (s *authService) GetCredentialByToken(token string) (*entity.Credential, error) {
	return s.authRepository.GetCredentialByToken(token)
}

func (s *authService) RemoveSingleSession(userID string, token string) error {
	return s.authRepository.RemoveSingleSession(userID, token)
}
