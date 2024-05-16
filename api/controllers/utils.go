package controllers

import (
	"encoding/json"
	"io"

	"github.com/labstack/echo/v4"
)

func getBodyInMap(c echo.Context) map[string]any {
	body, _ := io.ReadAll(c.Request().Body)
	var bodyMap map[string]any
	json.Unmarshal(body, &bodyMap)
	return bodyMap
}
