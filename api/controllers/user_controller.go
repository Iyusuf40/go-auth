package controllers

import (
	"fmt"
	"net/http"

	"github.com/Iyusuf40/go-auth/storage"
	"github.com/labstack/echo/v4"
)

var userStorage = storage.MakeUserStorage("")

func SaveUser(c echo.Context) error {
	body := getBodyInMap(c)
	userDesc, ok := body["data"].(map[string]any)
	response := map[string]string{}

	if !ok {
		response["error"] = "data payload is not decodeable into a map"
		return c.JSON(http.StatusBadRequest, response)
	}

	user := userStorage.BuildClient(userDesc)
	msg, success := userStorage.Save(user)

	if !success {
		response["error"] = msg
		return c.JSON(http.StatusBadRequest, response)
	} else {
		response["userId"] = msg
		return c.JSON(http.StatusCreated, response)
	}
}

func GetUser(c echo.Context) error {
	userId := c.Param("id")
	user, err := userStorage.Get(userId)
	response := map[string]string{}
	if err != nil {
		response["error"] = "user with id " + userId + " doesn't exist"
		return c.JSON(http.StatusNotFound, response)
	}
	return c.JSON(http.StatusOK, user)
}

func UpdateUser(c echo.Context) error {
	body := getBodyInMap(c)
	updateDesc, ok := body["data"].(map[string]any)
	response := map[string]string{}

	if !ok {
		response["error"] = "data payload is not decodeable into a map"
		return c.JSON(http.StatusBadRequest, response)
	}

	field, fieldOk := updateDesc["field"].(string)
	value := updateDesc["value"]

	if !fieldOk {
		response["error"] = "field part of updateDesc is not a string"
		return c.JSON(http.StatusBadRequest, response)
	}

	userId := c.Param("id")

	updated := userStorage.Update(userId, storage.UpdateDesc{Field: field,
		Value: value})

	if !updated {
		response["error"] = "update failed"
		return c.JSON(http.StatusBadRequest, response)
	}

	response["message"] = fmt.Sprintf("%s field of user succesfuly set to %s", field, value)
	return c.JSON(http.StatusOK, response)
}

func DeleteUser(c echo.Context) error {
	userId := c.Param("id")
	response := map[string]string{"message": "deleted"}

	userStorage.Delete(userId)
	return c.JSON(http.StatusOK, response)
}
