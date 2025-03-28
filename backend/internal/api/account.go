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

func GetAccount(c echo.Context) error {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*models.JwtUserClaims)

	db := c.Get("db").(*gorm.DB)
	var user models.User
	if err := db.Preload("Workouts.Exercises", "deleted = ?", false).Where("id = ?", claims.ID).First(&user).Error; err != nil {
		return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
	}

	uResponse := models.UserResponse{
		Username: user.Username,
		Email:    user.Email,

		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,

		Workouts: make([]models.WorkoutResponse, len(user.Workouts)),
	}

	for i, workout := range user.Workouts {
		wResponse := models.WorkoutResponse{
			ID: workout.ID,

			Name: workout.Name,
			Note: workout.Note,

			LastDone:  workout.LastDone,
			CreatedAt: workout.CreatedAt,
			UpdatedAt: workout.UpdatedAt,
			Exercises: make([]models.ExerciseResponse, len(workout.Exercises)),
		}

		for i, exercise := range workout.Exercises {
			eResponse := models.ExerciseResponse{
				ID:        exercise.ID,
				WorkoutID: exercise.WorkoutID,

				Name: exercise.Name,
				Note: exercise.Note,

				Sets: exercise.Sets,
				Reps: exercise.Reps,

				CreatedAt: workout.CreatedAt,
				UpdatedAt: workout.UpdatedAt,
			}
			wResponse.Exercises[i] = eResponse
		}
		uResponse.Workouts[i] = wResponse

	}

	return c.JSON(http.StatusOK, uResponse)
}

func GetToken(c echo.Context) error {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*models.JwtUserClaims)
	return c.JSON(http.StatusOK, echo.Map{"id": claims.ID, "username": claims.Username})
}

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
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
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
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	t, err := GenerateUserToken(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"token": t})
}

func DeleteAccount(c echo.Context) error {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*models.JwtUserClaims)

	db := c.Get("db").(*gorm.DB)
	var user models.User
	if err := db.Where("id = ?", claims.ID).First(&user).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("user_id = ?", user.ID).Delete(&models.ExerciseWeight{}).Error; err != nil {
			return err
		}

		if err := tx.Where("user_id = ?", user.ID).Delete(&models.ExerciseSession{}).Error; err != nil {
			return err
		}
		if err := tx.Where("user_id = ?", user.ID).Delete(&models.WorkoutSession{}).Error; err != nil {
			return err
		}
		if err := tx.Where("user_id = ?", user.ID).Delete(&models.Exercise{}).Error; err != nil {
			return err
		}
		if err := tx.Where("user_id = ?", user.ID).Delete(&models.Workout{}).Error; err != nil {
			return err
		}
		if err := tx.Delete(&user).Error; err != nil {
			return err
		}
		return nil // commit the transaction
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, "DELETED")
}
