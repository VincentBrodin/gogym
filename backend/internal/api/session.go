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

func StartSession(c echo.Context) error {
	workoutIDStr := c.Param("id")
	workoutID, err := strconv.ParseUint(workoutIDStr, 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid ID")
	}

	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*models.JwtUserClaims)

	db := c.Get("db").(*gorm.DB)

	var oldWorkoutSession models.WorkoutSession
	if err := db.Where("active = ? AND user_id = ?", true, claims.ID).First(&oldWorkoutSession).Error; err == nil {
		oldWorkoutSession.Active = false
		if err := db.Save(&oldWorkoutSession).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}

	var workout models.Workout
	if err := db.Preload("Exercises").Where("id = ? AND user_id = ?", workoutID, claims.ID).First(&workout).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	exerciseSessions := make([]models.ExerciseSession, len(workout.Exercises))
	for i, exercise := range workout.Exercises {
		exerciseSessions[i] = models.ExerciseSession{
			ExerciseID: exercise.ID,
			Exercise:   &exercise,

			Completed: false,
			Skiped:    false,
			Active:    false,

			SetsDone: 0,
		}
	}

	workoutSession := models.WorkoutSession{
		UserID: claims.ID,

		WorkoutID: workout.ID,
		Workout:   &workout,

		Active: true,

		StartedAt: time.Now().UTC(),
		EndedAt:   time.Now().UTC(),

		ExerciseSessions: exerciseSessions,
	}

	if err := db.Create(&workoutSession).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, workoutSession.CreateResponse())
}

func GetCurrentSession(c echo.Context) error {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*models.JwtUserClaims)

	db := c.Get("db").(*gorm.DB)

	var workoutSession models.WorkoutSession
	if err := db.Preload("Workout").Preload("ExerciseSessions").Preload("ExerciseSessions.Exercise").Where("active = ? AND user_id = ?", true, claims.ID).First(&workoutSession).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, workoutSession.CreateResponse())
}
