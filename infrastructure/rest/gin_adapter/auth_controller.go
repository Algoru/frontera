package ginadapter

import (
	"github.com/Algoru/frontera/domain/entity"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Algoru/frontera/configuration"
	"github.com/Algoru/frontera/domain/service"
	authrepository "github.com/Algoru/frontera/repository/auth_repository"
	"github.com/gin-gonic/gin"
)

func initAuthRoutes(group *gin.RouterGroup, as service.AuthService) {
	authPath := configuration.GetConfiguration().HTTP.AuthBasePath
	if authPath == "" {
		authPath = "/auth"
	}
	auth := group.Group(authPath)

	auth.POST("/login", loginController(as))
	auth.GET("/check", checkController(as))
	auth.POST("/refresh", refreshController(as))
	auth.DELETE("/logout", logoutController(as))
}

func loginController(as service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := new(authrepository.Auth)
		if err := c.BindJSON(auth); err != nil {
			log.Printf("unable to bind auth json: %s\n", err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		errors := auth.HasRequiredFields()
		if errors != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": errors})
			return
		}

		credential, err := as.Login(auth)
		if err != nil {
			log.Printf("unable to login user with email \"%s\": %s\n", auth.Email, err.Error())
			statusCode := entity.GetStatusCodeForError(err)
			c.AbortWithStatusJSON(statusCode, gin.H{"error": err.Error()})
			return
		}

		if !configuration.GetConfiguration().User.AllowMultipleSessions {
			if err := as.RemoveUserSessions(credential.UserID); err != nil {
				log.Printf("unable to remove user sessions: %s\n", err)
			}
		}

		if err = as.AddUserSession(credential); err != nil {
			log.Printf("unable to save user session: %s\n", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, credential)
	}
}

func logoutController(as service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader, err := getTokenFromAuthHeader(c)
		if err != nil {
			statusCode := entity.GetStatusCodeForError(err)
			c.AbortWithStatusJSON(statusCode, gin.H{"error": err.Error()})
			return
		}

		flush, err := strconv.ParseBool(c.Query("flush"))
		if err != nil {
			flush = false // Just to be explicit
		}

		if err := as.Logout(authHeader, flush); err != nil {
			log.Printf("unable to perform logout: %s\n", err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Status(http.StatusOK)
	}
}

func checkController(as service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader, err := getTokenFromAuthHeader(c)
		if err != nil {
			statusCode := entity.GetStatusCodeForError(err)
			c.AbortWithStatusJSON(statusCode, gin.H{"error": err.Error()})
			return
		}

		credential, err := as.GetCredentialByToken(authHeader)
		if err != nil {
			log.Printf("unable to get credential by token: %s\n", err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenDuration := (time.Duration(configuration.GetConfiguration().Security.TokenLifetime) * time.Hour).Hours()
		hoursFromNow := time.Now().Sub(credential.ExpiresAt).Hours()
		if hoursFromNow >= tokenDuration {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Status(http.StatusOK)
	}
}

func refreshController(as service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader, err := getTokenFromAuthHeader(c)
		if err != nil {
			statusCode := entity.GetStatusCodeForError(err)
			c.AbortWithStatusJSON(statusCode, gin.H{"error": err.Error()})
			return
		}

		credential, err := as.GetCredentialByToken(authHeader)
		if err != nil {
			log.Printf("unable to get credential by token: %s\n", err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		refreshed, err := as.RefreshCredential(credential)
		if err != nil {
			log.Printf("unable to refresh credential: %s\n", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, refreshed)
	}
}

func getTokenFromAuthHeader(c *gin.Context) (string, error) {
	headerValue := c.GetHeader("Authorization")
	if headerValue == "" {
		return "", entity.ErrInvalidAuthHeader
	}

	parts := strings.Split(headerValue, " ")
	if len(parts) != 2 {
		return "", entity.ErrInvalidAuthHeader
	}

	return parts[1], nil
}