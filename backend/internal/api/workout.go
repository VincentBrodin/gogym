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

type AddWorkoutForm struct {
	Name string `json:"name"`
	Note string `json:"note"`
}

func GetAllWorkouts(c echo.Context) error {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*models.JwtUserClaims)

	db := c.Get("db").(*gorm.DB)

	workouts := make([]*models.Workout, 0)
	if err := db.Preload("Exercises", "deleted = ?", false).Where("user_id = ? AND deleted = ?", claims.ID, false).Find(&workouts).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	wResponses := make([]models.WorkoutResponse, len(workouts))

	for i, workout := range workouts {
		wResponses[i] = workout.CreateResponse()
	}

	return c.JSON(http.StatusOK, wResponses)
}

func GetWorkout(c echo.Context) error {
	workoutIDStr := c.Param("id")
	workoutID, err := strconv.ParseUint(workoutIDStr, 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid ID")
	}

	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*models.JwtUserClaims)

	db := c.Get("db").(*gorm.DB)

	var workout models.Workout
	if err := db.Preload("Exercises", "deleted = ?", false).Where("id = ? AND user_id = ? AND deleted = ?", workoutID, claims.ID, false).First(&workout).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, workout.CreateResponse())
}

func AddWorkout(c echo.Context) error {
	var form AddWorkoutForm
	if err := c.Bind(&form); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*models.JwtUserClaims)

	workout := models.Workout{
		UserID: claims.ID,
		Name:   form.Name,
		Note:   form.Note,

		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		LastDone:  time.Now().UTC(),
	}

	db := c.Get("db").(*gorm.DB)
	if err := db.Create(&workout).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, workout.CreateResponse())
}

func EditWorkout(c echo.Context) error {
	var form AddWorkoutForm
	if err := c.Bind(&form); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	workoutIDStr := c.Param("id")
	workoutID, err := strconv.ParseUint(workoutIDStr, 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid ID")
	}

	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*models.JwtUserClaims)

	db := c.Get("db").(*gorm.DB)

	var workout models.Workout
	if err := db.Where("id = ? AND user_id = ? AND deleted = ?", workoutID, claims.ID, false).First(&workout).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	workout.Name = form.Name
	workout.Note = form.Note
	workout.UpdatedAt = time.Now().UTC()

	if err := db.Save(&workout).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, workout.CreateResponse())
}

func DeleteWorkout(c echo.Context) error {
	workoutIDStr := c.Param("id")
	workoutID, err := strconv.ParseUint(workoutIDStr, 10, 32)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid ID")
	}

	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*models.JwtUserClaims)

	db := c.Get("db").(*gorm.DB)

	var workout models.Workout
	if err := db.Preload("Exercises", "deleted = ?", false).Where("id = ? AND user_id = ? AND deleted = ?", workoutID, claims.ID, false).First(&workout).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	// Update order and save
	err = db.Transaction(func(tx *gorm.DB) error {
		workout.Deleted = true
		workout.DeletedAt = time.Now().UTC()
		if err := tx.Save(&workout).Error; err != nil {
			return err
		}
		for i := range workout.Exercises {
			workout.Exercises[i].Deleted = true
			workout.Exercises[i].DeletedAt = time.Now().UTC()
			if err := tx.Save(&workout.Exercises[i]).Error; err != nil {
				return err // transaction will be rolled back
			}
		}
		return nil // commit the transaction
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, "deleted")
}
