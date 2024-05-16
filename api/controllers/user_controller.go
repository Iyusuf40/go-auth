package controllers

import (
	"net/http"

	"github.com/Iyusuf40/go-auth/storage"
	"github.com/labstack/echo/v4"
)

var userStorage = storage.MakeUserStorage("")

func SaveUser(c echo.Context) error {
	body := getBodyInMap(c)
	userDesc := body["data"].(map[string]any)
	user := userStorage.BuildClient(userDesc)

	msg, success := userStorage.Save(user)

	response := map[string]string{}

	if !success {
		response["error"] = msg
		return c.JSON(http.StatusBadRequest, response)
	} else {
		response["userId"] = msg
		return c.JSON(http.StatusCreated, response)
	}
}
