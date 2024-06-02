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

var AUTH_HANDLER = MakeAuthHandler("", "", "")

func Login(c echo.Context) error {
	body := controllers.GetBodyInMap(c)
	userDesc, ok := body["data"].(map[string]any)
	response := map[string]string{}

	if !ok {
		response["error"] = "data payload is not decodeable into a map"
		return c.JSON(http.StatusBadRequest, response)
	}

	email, email_ok := userDesc["email"].(string)

	if !email_ok {
		response["error"] = "email required to login"
		return c.JSON(http.StatusBadRequest, response)
	}

	password, password_ok := userDesc["password"].(string)

	if !password_ok {
		response["error"] = "password required to login"
		return c.JSON(http.StatusBadRequest, response)
	}

	sessionId := AUTH_HANDLER.HandleLogin(email, password)

	if sessionId == "" {
		response["error"] = "failed to login"
		return c.JSON(http.StatusBadRequest, response)
	}
	response["sessionId"] = sessionId
	return c.JSON(http.StatusOK, response)
}

func Logout(c echo.Context) error {
	body := controllers.GetBodyInMap(c)
	userDesc, ok := body["data"].(map[string]any)
	response := map[string]string{}

	if !ok {
		response["error"] = "data payload is not decodeable into a map"
		return c.JSON(http.StatusBadRequest, response)
	}

	sessionId, sessionId_ok := userDesc["sessionId"].(string)

	if !sessionId_ok {
		response["error"] = "sessionId required to logout"
		return c.JSON(http.StatusBadRequest, response)
	}

	AUTH_HANDLER.HandleLogout(sessionId)
	response["message"] = "logged out"
	return c.JSON(http.StatusOK, response)
}

func IsLoggedIn(c echo.Context) error {
	body := controllers.GetBodyInMap(c)
	userDesc, ok := body["data"].(map[string]any)
	response := map[string]any{}

	if !ok {
		response["error"] = "data payload is not decodeable into a map"
		return c.JSON(http.StatusBadRequest, response)
	}

	sessionId, sessionId_ok := userDesc["sessionId"].(string)

	if !sessionId_ok {
		response["error"] = "sessionId required to check if user is logged in"
		return c.JSON(http.StatusBadRequest, response)
	}

	isLoggedIn := AUTH_HANDLER.IsLoggedIn(sessionId)
	response["isLoggedIn"] = isLoggedIn

	return c.JSON(http.StatusOK, response)
}
