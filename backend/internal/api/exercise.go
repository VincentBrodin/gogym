package api

import (
	"backend/internal/models"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type AddExerciseForm struct {
	WorkoutID uint    `form:"workout_id"`
	Name      string  `form:"name"`
	Note      *string `form:"note"`

	Sets int `form:"sets"`
	Reps int `form:"reps"`
}

func GetExercise(c echo.Context) error {
	exerciseIDStr := c.Param("id")
	exerciseID, err := strconv.ParseUint(exerciseIDStr, 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid ID")
	}

	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*models.JwtUserClaims)

	db := c.Get("db").(*gorm.DB)

	var exercise models.Exercise
	if err := db.Where("id = ? AND user_id = ?", exerciseID, claims.ID).First(&exercise).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, exercise.CreateResponse())
}

func AddExercise(c echo.Context) error {
	var form AddExerciseForm
	if err := c.Bind(&form); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*models.JwtUserClaims)

	exercise := models.Exercise{
		UserID:    claims.ID,
		WorkoutID: form.WorkoutID,

		Name: form.Name,
		Note: form.Note,

		Sets: form.Sets,
		Reps: form.Reps,

		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	db := c.Get("db").(*gorm.DB)
	if err := db.Create(&exercise).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"id": exercise.ID})
}

func EditExercise(c echo.Context) error {
	var form AddExerciseForm
	if err := c.Bind(&form); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	exerciseIDStr := c.Param("id")
	exerciseID, err := strconv.ParseUint(exerciseIDStr, 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid ID")
	}

	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*models.JwtUserClaims)

	db := c.Get("db").(*gorm.DB)

	var exercise models.Exercise
	if err := db.Where("id = ? AND user_id = ?", exerciseID, claims.ID).First(&exercise).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	// exercise.WorkoutID = form.WorkoutID

	exercise.Name = form.Name
	exercise.Note = form.Note
	exercise.Sets = form.Sets
	exercise.Reps = form.Reps

	exercise.UpdatedAt = time.Now().UTC()

	if err := db.Save(&exercise).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, "updated")
}

func DeleteExercise(c echo.Context) error {
	exerciseIDStr := c.Param("id")
	exerciseID, err := strconv.ParseUint(exerciseIDStr, 10, 32)

	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid ID")
	}

	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*models.JwtUserClaims)

	db := c.Get("db").(*gorm.DB)
	if err := db.Where("id = ? AND user_id = ?", exerciseID, claims.ID).Delete(&models.Exercise{}).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, "deleted")
}
