package ginadapter

import (
	"fmt"

	"github.com/Algoru/frontera/configuration"

	"github.com/Algoru/frontera/domain/service"
	"github.com/gin-gonic/gin"
)

// GinAdapter
type GinAdapter struct {
	Port        uint16
	router      *gin.Engine
	userService service.UserService
	authService service.AuthService
}

// NewGinAdapter
func NewGinAdapter(port uint16) GinAdapter {
	if !configuration.GetConfiguration().Debug {
		gin.SetMode(gin.ReleaseMode)
	}

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

// InitRoutes
func (ga *GinAdapter) InitRoutes() {
	root := ga.router.Group("/")
	initUserRoutes(root, ga.userService)
	initAuthRoutes(root, ga.authService)
}

// SetUserService
func (ga *GinAdapter) SetUserService(us service.UserService) {
	ga.userService = us
}

// SetAuthService
func (ga *GinAdapter) SetAuthService(as service.AuthService) {
	ga.authService = as
}
