package api

import (
	"backend/internal/models"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type AddExerciseForm struct {
	WorkoutID uint    `json:"workout_id"`
	Name      string  `json:"name"`
	Note      *string `json:"note"`

	Order int `json:"order"`

	Sets int `json:"sets"`
	Reps int `json:"reps"`
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
	// Users could add exercises to workouts that they do not own!!
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

		Order: form.Order,

		Sets: form.Sets,
		Reps: form.Reps,

		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	db := c.Get("db").(*gorm.DB)
	if err := db.Create(&exercise).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Send back exercise
	return c.JSON(http.StatusOK, exercise.CreateResponse())
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

	var exercise models.Exercise
	if err := db.Where("id = ? AND user_id = ?", exerciseID, claims.ID).First(&exercise).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	if err := db.Delete(&exercise).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	var exercises []models.Exercise
	if err := db.Where("workout_id = ? AND user_id = ?", exercise.WorkoutID, claims.ID).Find(&exercises).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	sort.Slice(exercises, func(i, j int) bool {
		return exercises[i].Order < exercises[j].Order
	})

	for index, ex := range exercises {
		ex.Order = index
		if err := db.Save(&ex).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}

	return c.JSON(http.StatusOK, "deleted")
}
