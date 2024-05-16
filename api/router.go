package api

import (
	"github.com/Iyusuf40/go-auth/api/controllers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func Serve() {
	e := echo.New()
	e.Use(middleware.Recover())
	g := e.Group("/api")
	g.POST("/users", controllers.SaveUser)
	e.Logger.Fatal(e.Start(":8080"))
}
