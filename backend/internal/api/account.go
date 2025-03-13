package api

import (
	"backend/internal/models"
	"net/http"
	"net/mail"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func EditAccount(c echo.Context) error {
	var form RegistrationForm
	if err := c.Bind(&form); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*models.JwtUserClaims)

	db := c.Get("db").(*gorm.DB)
	var user models.User
	if err := db.Where("id = ?", claims.ID).First(&user).Error; err != nil {
		return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
	}

	uname := len(form.Username)
	if 4 < uname && uname < 256 {
		user.Username = form.Username
	}

	if _, err := mail.ParseAddress(form.Email); err == nil {
		user.Email = form.Email
	}

	user.UpdatedAt = time.Now().UTC()

	if err := db.Save(&user).Error; err != nil {
		return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
	}

	t, err := GenerateUserToken(user)
	if err != nil {
		return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"token": t})
}
