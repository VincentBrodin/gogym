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
	if err := db.Preload("Exercises", "deleted = ?", false).Where("id = ? AND user_id = ?", workoutID, claims.ID).First(&workout).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	workout.LastDone = time.Now().UTC()

	if err := db.Save(&workout).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	exerciseSessions := make([]models.ExerciseSession, len(workout.Exercises))
	for i, exercise := range workout.Exercises {
		exerciseSessions[i] = models.ExerciseSession{
			UserID: claims.ID,
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

func EditSession(c echo.Context) error {
	sessionIDStr := c.Param("id")
	sessionID, err := strconv.ParseUint(sessionIDStr, 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid ID")
	}

	var form models.WorkoutSessionResponse
	if err := c.Bind(&form); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*models.JwtUserClaims)

	db := c.Get("db").(*gorm.DB)

	var workoutSession models.WorkoutSession
	if err := db.Preload("Workout", "deleted = ?", false).Preload("ExerciseSessions").Preload("ExerciseSessions.Exercise", "deleted = ?", false).Where("id = ? AND user_id = ?", sessionID, claims.ID).First(&workoutSession).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	workoutSession.Active = form.Active
	workoutSession.EndedAt = time.Now().UTC()

	// Create a map to sync incoming with db
	exerciseMap := make(map[uint]*models.ExerciseSession)
	for i := range workoutSession.ExerciseSessions {
		exerciseMap[workoutSession.ExerciseSessions[i].ID] = &workoutSession.ExerciseSessions[i]
	}

	for _, newExercise := range form.ExerciseSessions {
		if oldExercise, ok := exerciseMap[newExercise.ID]; ok {
			oldExercise.Completed = newExercise.Completed
			oldExercise.Skiped = newExercise.Skiped
			oldExercise.Active = newExercise.Active
			oldExercise.SetsDone = newExercise.SetsDone
		}
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&workoutSession).Error; err != nil {
			return err
		}
		for i := range workoutSession.ExerciseSessions {
			if err := tx.Save(&workoutSession.ExerciseSessions[i]).Error; err != nil {
				return err // transaction will be rolled back
			}
		}
		return nil // commit the transaction
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, workoutSession.CreateResponse())
}

func GetCurrentSession(c echo.Context) error {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*models.JwtUserClaims)

	db := c.Get("db").(*gorm.DB)

	var workoutSession models.WorkoutSession
	if err := db.Preload("Workout", "deleted = ?", false).Preload("ExerciseSessions").Preload("ExerciseSessions.Exercise", "deleted = ?", false).Where("active = ? AND user_id = ?", true, claims.ID).First(&workoutSession).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, workoutSession.CreateResponse())
}
