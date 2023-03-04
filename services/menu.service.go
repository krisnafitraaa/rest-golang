package services

import (
	"net/http"

	"github.com/fundraising/rest-api/repositories"
	"github.com/labstack/echo"
)

func GetAllMenus(c echo.Context) error {
	response, err := repositories.FetchAllMenus()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, response)
}

func GetMenu(c echo.Context) error {
	response, err := repositories.GetMenu()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, response)
}

func StoreMenu(c echo.Context) error {
	params := map[string]interface{}{
		"menu":      c.FormValue("menu"),
		"icon":      c.FormValue("icon"),
		"position":  c.FormValue("position"),
		"menu_for":  c.FormValue("menu_for"),
		"menu_url":  c.FormValue("menu_url"),
		"is_active": c.FormValue("is_active"),
	}

	response, err := repositories.StoreMenu(params)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, response)
}
