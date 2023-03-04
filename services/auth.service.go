package services

import (
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/fundraising/rest-api/helpers"
	"github.com/fundraising/rest-api/repositories"
	"github.com/google/uuid"
	"github.com/labstack/echo"
)

func DoRegister(c echo.Context) error {
	user_id := strings.Split(uuid.New().String(), "-")[0]
	volunteer_id := strings.Split(uuid.New().String(), "-")[0]
	password, _ := helpers.HashPassword(c.FormValue("password"))

	params := map[string]interface{}{
		"user_id":       user_id,
		"volunteer_id":  volunteer_id,
		"role_id":       c.FormValue("role_id"),
		"fullname":      c.FormValue("fullname"),
		"email":         c.FormValue("email"),
		"password":      password,
		"photo":         c.FormValue("photo"),
		"is_loggedin":   0,
		"provider":      c.FormValue("provider"),
		"is_active":     c.FormValue("is_active"),
		"verified_user": 0,
		"created_by":    c.FormValue("created_by"),
	}

	response, err := repositories.RegisterNewUser(params)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, response)
}

func DoLogin(c echo.Context) error {
	email := c.FormValue("email")
	password := c.FormValue("password")

	params := map[string]string{
		"email":    email,
		"password": password,
	}

	response, err := repositories.CheckLogin(params)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, response)
	}

	//generate token from jwt
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = email
	claims["expired"] = time.Now().Add(time.Hour * 2).Unix()

	t, err := token.SignedString([]byte("nothingknown"))

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  200,
		"message": "login success",
		"valid":   true,
		"token":   t,
	})
}
