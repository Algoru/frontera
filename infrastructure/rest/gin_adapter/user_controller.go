package ginadapter

import (
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Algoru/frontera/configuration"
	"github.com/Algoru/frontera/domain/service"
	userrepository "github.com/Algoru/frontera/repository/user_repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func initUserRoutes(group *gin.RouterGroup, us service.UserService) {
	usersPath := configuration.GetConfiguration().HTTP.UserBasePath
	if usersPath == "" {
		usersPath = "/users"
	}
	users := group.Group(usersPath)

	users.POST("/", createUserController(us))
	users.GET("/:id", getUserController(us))
	users.PUT("/:id", updateUserController(us))
	users.DELETE("/:id", deleteUserController(us))
}

func createUserController(us service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := new(userrepository.User)
		if err := c.BindJSON(user); err != nil {
			log.Printf("unable to bind user JSON: %s\n", err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user.Sanitize()
		if errors := us.HasRequiredFields(user); errors != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": errors})
			return
		}

		_, err := us.GetUserByEmail(user.Email)
		if err == nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "email already being used"})
			return
		} else if err != nil && err != mongo.ErrNoDocuments {
			log.Printf("unable to check if email is being used: %s\n", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"errors": "unable to check if email is in use"})
			return
		}

		created, err := us.CreateUser(user)
		if err != nil {
			log.Printf("unable to save user: %s\n", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		created.RemoveSensible()

		c.JSON(http.StatusCreated, created)
	}
}

func getUserController(us service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			log.Printf("unable to parse \"%s\" as uuid: %s\n", c.Param("id"), err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := us.GetUser(userID)
		if err != nil {
			// TODO (@Algoru): Identify if it was an internal error or if user couldn't be found.
			log.Printf("unable to get user with id \"%s\": %s\n", c.Param("id"), err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

func updateUserController(us service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.String(http.StatusCreated, "user updated")
	}
}

func deleteUserController(us service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			log.Printf("unable to parse \"%s\" as uuid: %s\n", c.Param("id"), err)
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		deletedUser, err := us.DeleteUser(userID)
		if err != nil {
			// TODO (@Algoru): Identify if it was an internal error or if user couldn't be found.
			log.Printf("unable to delete user with uuid \"%s\": %s\n", c.Param("id"), err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, deletedUser)
	}
}
