package api

import (
	"backend/internal/models"
	"fmt"
	"net/http"
	"sort"

	// "sort"
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

	now := time.Now().UTC()
	workout.LastDone = now

	workoutSession := models.WorkoutSession{
		UserID: claims.ID,

		WorkoutID: workout.ID,

		Active: true,

		StartedAt:        now,
		EndedAt:          now,
		ExerciseSessions: []models.ExerciseSession{},
	}

	err = db.Transaction(func(tx *gorm.DB) error {

		if err := db.Save(&workout).Error; err != nil {
			return err
		}

		if err := db.Create(&workoutSession).Error; err != nil {
			return err
		}

		for _, exercise := range workout.Exercises {
			exerciseSession := models.ExerciseSession{
				UserID:           claims.ID,
				ExerciseID:       exercise.ID,
				WorkoutSessionID: workoutSession.ID,

				Completed: false,
				Skiped:    false,
				Active:    false,

				SetsDone:        0,
				ExerciseWeights: []models.ExerciseWeight{},
			}
			if err := db.Create(&exerciseSession).Error; err != nil {
				return err
			}

			fmt.Println(exerciseSession.ID)
			var weight float64 = 0
			var lastWeight models.ExerciseWeight
			if err := db.Joins("JOIN exercises ON exercises.id = exercise_weights.exercise_id").
				Where("exercises.name = ? AND exercise_weights.user_id = ?", exercise.Name, claims.ID).
				Order("exercise_weights.id DESC").
				First(&lastWeight).Error; err == nil {
				weight = lastWeight.Weight
			}
			for j := range exercise.Sets {
				exerciseWeight := models.ExerciseWeight{
					UserID:            claims.ID,
					ExerciseID:        exercise.ID,
					ExerciseSessionID: exerciseSession.ID,
					Set:               j + 1,
					Weight:            weight,
				}
				if err := db.Create(&exerciseWeight).Error; err != nil {
					return err
				}
			}

		}
		return nil // commit the transaction
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	var workoutSessionTotal models.WorkoutSession
	if err := db.Preload("Workout", "deleted = ?", false).Preload("ExerciseSessions").Preload("ExerciseSessions.Exercise", "deleted = ?", false).Preload("ExerciseSessions.ExerciseWeights").Where("id = ? AND user_id = ?", workoutSession.ID, claims.ID).First(&workoutSessionTotal).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, workoutSessionTotal.CreateResponse())
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
	if err := db.Preload("Workout", "deleted = ?", false).Preload("ExerciseSessions").Preload("ExerciseSessions.Exercise", "deleted = ?", false).Preload("ExerciseSessions.ExerciseWeights").Where("id = ? AND user_id = ?", sessionID, claims.ID).First(&workoutSession).Error; err != nil {
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
			sort.Slice(newExercise.ExerciseWeights, func(i, j int) bool {
				return oldExercise.ExerciseWeights[i].Set < oldExercise.ExerciseWeights[j].Set
			})

			sort.Slice(oldExercise.ExerciseWeights, func(i, j int) bool {
				return oldExercise.ExerciseWeights[i].Set < oldExercise.ExerciseWeights[j].Set
			})

			for i := range newExercise.ExerciseWeights {
				oldExercise.ExerciseWeights[i].Weight = newExercise.ExerciseWeights[i].Weight
			}
		}
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&workoutSession).Error; err != nil {
			return err
		}
		for i := range workoutSession.ExerciseSessions {
			for j := range workoutSession.ExerciseSessions[i].ExerciseWeights {
				if err := tx.Save(&workoutSession.ExerciseSessions[i].ExerciseWeights[j]).Error; err != nil {
					return err // transaction will be rolled back
				}
			}

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

func GetAllSessions(c echo.Context) error {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*models.JwtUserClaims)

	db := c.Get("db").(*gorm.DB)
	var workoutSessions []models.WorkoutSession
	if err := db.Preload("Workout").Preload("ExerciseSessions").Preload("ExerciseSessions.Exercise").Preload("ExerciseSessions.ExerciseWeights").Where("active = ? AND user_id = ?", false, claims.ID).Find(&workoutSessions).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	response := make([]models.WorkoutSessionResponse, len(workoutSessions))
	for i, session := range workoutSessions {
		if session.Workout == nil {
			fmt.Printf("Session has no workout %d\n", session.ID)
			continue
		}
		response[i] = session.CreateResponse()
	}

	return c.JSON(http.StatusOK, response)
}

func GetSession(c echo.Context) error {
	sessionIDStr := c.Param("id")
	sessionID, err := strconv.ParseUint(sessionIDStr, 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid ID")
	}

	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*models.JwtUserClaims)

	db := c.Get("db").(*gorm.DB)
	var workoutSession models.WorkoutSession
	if err := db.Preload("Workout").Preload("ExerciseSessions").Preload("ExerciseSessions.Exercise").Preload("ExerciseSessions.ExerciseWeights").Where("id = ? AND user_id = ?", sessionID, claims.ID).First(&workoutSession).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, workoutSession.CreateResponse())
}

func GetCurrentSession(c echo.Context) error {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*models.JwtUserClaims)

	db := c.Get("db").(*gorm.DB)

	var workoutSession models.WorkoutSession
	if err := db.Preload("Workout", "deleted = ?", false).Preload("ExerciseSessions").Preload("ExerciseSessions.Exercise", "deleted = ?", false).Preload("ExerciseSessions.ExerciseWeights").Where("active = ? AND user_id = ?", true, claims.ID).First(&workoutSession).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, workoutSession.CreateResponse())
}

func DeleteSession(c echo.Context) error {
	sessionIDStr := c.Param("id")
	sessionID, err := strconv.ParseUint(sessionIDStr, 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid ID")
	}

	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*models.JwtUserClaims)

	db := c.Get("db").(*gorm.DB)

	var workoutSession models.WorkoutSession
	if err := db.Preload("ExerciseSessions").Preload("ExerciseSessions.ExerciseWeights").Where("id = ? AND user_id = ?", sessionID, claims.ID).First(&workoutSession).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	err = db.Transaction(func(tx *gorm.DB) error {

		for i := range workoutSession.ExerciseSessions {
			for j := range workoutSession.ExerciseSessions[i].ExerciseWeights {
				if err := tx.Delete(&workoutSession.ExerciseSessions[i].ExerciseWeights[j]).Error; err != nil {
					return err // transaction will be rolled back
				}
			}
			if err := tx.Delete(&workoutSession.ExerciseSessions[i]).Error; err != nil {
				return err // transaction will be rolled back
			}
		}

		if err := tx.Delete(&workoutSession).Error; err != nil {
			return err
		}

		return nil // commit the transaction
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, "DELETED")
}
