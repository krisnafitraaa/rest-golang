package services

import (
	"net/http"

	"github.com/fundraising/rest-api/repositories"
	"github.com/labstack/echo"
)

func GetRole(c echo.Context) error {
	response, err := repositories.GetRole()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  false,
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, response)
}

func RoleDetail(c echo.Context) error {
	role_id := c.Param("roleid")

	response, err := repositories.RoleDetail(role_id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  false,
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, response)
}

func StoreRole(c echo.Context) error {
	params := map[string]interface{}{
		"role_id":    c.FormValue("role_id"),
		"role":       c.FormValue("role"),
		"role_group": c.FormValue("role_group"),
		"added_info": c.FormValue("added_info"),
	}

	response, err := repositories.StoreRole(params)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  false,
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, response)
}

func UpdateRole(c echo.Context) error {
	params := map[string]interface{}{
		"key":        c.FormValue("key"),
		"role_id":    c.FormValue("role_id"),
		"role":       c.FormValue("role"),
		"role_group": c.FormValue("role_group"),
		"added_info": c.FormValue("added_info"),
	}

	response, err := repositories.UpdateRole(params)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"status":  false,
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, response)
}

func DeleteRole(c echo.Context) error {
	role_id := c.FormValue("role_id")

	response, err := repositories.DeleteRole(role_id)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": err.Error(),
			"status":  false,
		})
	}

	return c.JSON(http.StatusOK, response)
}
