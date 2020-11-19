package ginadapter

import (
	"log"
	"net/http"
	"strconv"
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

		credential, err := as.Login(auth)
		if err != nil {
			log.Printf("unable to login user with email \"%s\": %s\n", auth.Email, err.Error())
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
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
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		flush, err := strconv.ParseBool(c.Query("flush"))
		if err != nil {
			flush = false // Just to be explicit
		}

		if err := as.Logout(token, flush); err != nil {
			log.Printf("unable to perform logout: %s\n", err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Status(http.StatusOK)
	}
}

func checkController(as service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		credential, err := as.GetCredentialByToken(token)
		if err != nil {
			log.Printf("unable to get credential by token: %s\n", err)
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenDuration := time.Duration(configuration.GetConfiguration().Security.TokenLifetime) * time.Hour
		if time.Now().Sub(credential.ExpiresAt) >= tokenDuration {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.Status(http.StatusOK)
	}
}

func refreshController(as service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
