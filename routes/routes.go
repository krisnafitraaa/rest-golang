package routes

import (
	"net/http"

	"github.com/fundraising/rest-api/middleware"
	"github.com/fundraising/rest-api/services"
	"github.com/labstack/echo"
)

func DefineRoutes() *echo.Echo {
	route := echo.New()

	route.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome to kamipeduli web service")
	})

	//get
	route.GET("/menus", services.GetAllMenus)
	route.GET("/menu", services.GetMenu)
	route.GET("/users", services.GetUser, middleware.IsAuthenticated)
	route.GET("/roles", services.GetRole, middleware.IsAuthenticated)
	route.GET("/role/detail/:roleid", services.RoleDetail, middleware.IsAuthenticated)

	//post
	route.POST("/menu", services.StoreMenu, middleware.IsAuthenticated)
	route.POST("/role", services.StoreRole, middleware.IsAuthenticated)
	route.POST("/register", services.DoRegister)
	route.POST("/login", services.DoLogin)

	//put
	route.PUT("/role", services.UpdateRole, middleware.IsAuthenticated)

	//delete
	route.DELETE("/role", services.DeleteRole, middleware.IsAuthenticated)

	return route
}
