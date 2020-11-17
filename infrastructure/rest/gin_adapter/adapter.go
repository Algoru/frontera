package ginadapter

import (
	"fmt"

	"github.com/Algoru/frontera/domain/service"
	"github.com/gin-gonic/gin"
)

// GinAdapter
type GinAdapter struct {
	Port        uint16
	router      *gin.Engine
	userService service.UserService
}

func NewGinAdapter(port uint16) GinAdapter {
	r := gin.Default()

	return GinAdapter{
		Port:   port,
		router: r,
	}
}

// Start
func (ga *GinAdapter) Start() error {
	addr := fmt.Sprintf(":%d", ga.Port)
	return ga.router.Run(addr)
}

func (ga *GinAdapter) SetUserService(us service.UserService) {
	ga.userService = us
}

func (ga *GinAdapter) InitRoutes() {
	root := ga.router.Group("/")
	InitUserRoutes(root, ga.userService)
}
