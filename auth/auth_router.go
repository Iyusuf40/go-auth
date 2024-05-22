package auth

import (
	"net/http"

	"github.com/Iyusuf40/go-auth/api/controllers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func ServeAUTH() {
	e := echo.New()
	e.Use(middleware.Recover())

	g := e.Group("/auth")
	g.POST("/login", Login)
	g.GET("/logout", Logout)
	g.PUT("/isloggedin", IsLoggedIn)

	e.Logger.Fatal(e.Start(":8081"))
}

func Login(c echo.Context) error {
	body := controllers.GetBodyInMap(c)
	userDesc, ok := body["data"].(map[string]any)
	response := map[string]string{}

	if !ok {
		response["error"] = "data payload is not decodeable into a map"
		return c.JSON(http.StatusBadRequest, response)
	}

	if userDesc != nil && response != nil {

	}
	return nil
}

func Logout(c echo.Context) error {
	body := controllers.GetBodyInMap(c)
	userDesc, ok := body["data"].(map[string]any)
	response := map[string]string{}

	if !ok {
		response["error"] = "data payload is not decodeable into a map"
		return c.JSON(http.StatusBadRequest, response)
	}

	if userDesc != nil && response != nil {

	}
	return nil
}

func IsLoggedIn(c echo.Context) error {
	body := controllers.GetBodyInMap(c)
	userDesc, ok := body["data"].(map[string]any)
	response := map[string]string{}

	if !ok {
		response["error"] = "data payload is not decodeable into a map"
		return c.JSON(http.StatusBadRequest, response)
	}

	if userDesc != nil && response != nil {

	}
	return nil
}
