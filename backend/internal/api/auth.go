package api

import (
	"net/http"
	"net/mail"
	"os"
	"time"

	"backend/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type RegistrationForm struct {
	Username string `form:"uname"`
	Email    string `form:"email"`
	Password string `form:"pswd"`
}

type LoginForm struct {
	Username string `form:"uname"`
	Password string `form:"pswd"`
}

func (form *RegistrationForm) Validate() map[string]string {
	errors := make(map[string]string)
	uname := len(form.Username)
	if 4 > uname || uname > 256 {
		errors["uname"] = "Too big or small (min 4, max 256)"
	}

	if _, err := mail.ParseAddress(form.Email); err != nil {
		errors["email"] = "Invalid email adress"
	}

	if len(errors) == 0 {
		return nil
	}

	return errors
}

func Login(c echo.Context) error {
	var form LoginForm
	if err := c.Bind(&form); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	db := c.Get("db").(*gorm.DB)
	var user models.User
	if err := db.Where("username = ?", form.Username).First(&user).Error; err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(form.Password)); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Username or password is worng"})
	}

	t, err := GenerateUserToken(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"token": t})
}

func Register(c echo.Context) error {
	var form RegistrationForm
	if err := c.Bind(&form); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	if val := form.Validate(); val != nil {
		return c.JSON(http.StatusBadRequest, val)
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte(form.Password), 14)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Could not hash password"})
	}

	user := models.User{
		Username:  form.Username,
		Email:     form.Email,
		Password:  string(bytes),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	db := c.Get("db").(*gorm.DB)
	if err := db.Create(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	t, err := GenerateUserToken(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"token": t})
}

func GenerateUserToken(user models.User) (string, error) {
	// Create token
	claims := &models.JwtUserClaims{
		ID:       user.ID,
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 72)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}
	return t, nil
}
