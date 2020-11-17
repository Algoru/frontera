package ginadapter

import (
	"log"
	"net/http"

	"github.com/Algoru/frontera/domain/service"
	userrepository "github.com/Algoru/frontera/repository/user_repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func InitUserRoutes(group *gin.RouterGroup, us service.UserService) {
	users := group.Group("/users")

	log.Println("us when init:", us)
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

		created, err := us.CreateUser(user)
		if err != nil {
			log.Printf("unable to save user: %s\n", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

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
