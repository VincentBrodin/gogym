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
	Name string  `form:"name"`
	Note *string `form:"note"`
}

func GetAllWorkouts(c echo.Context) error {
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*models.JwtUserClaims)

	db := c.Get("db").(*gorm.DB)

	workouts := make([]*models.Workout, 0)
	if err := db.Preload("Exercises").Where("user_id = ?", claims.ID).Find(&workouts).Error; err != nil {
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
	if err := db.Preload("Exercises").Where("id = ? AND user_id = ?", workoutID, claims.ID).First(&workout).Error; err != nil {
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

	return c.JSON(http.StatusOK, echo.Map{"id": workout.ID})
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
	if err := db.Where("id = ? AND user_id = ?", workoutID, claims.ID).First(&workout).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	workout.Name = form.Name
	workout.Note = form.Note
	workout.UpdatedAt = time.Now().UTC()

	if err := db.Save(&workout).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, "updated")
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

	if err := db.Where("workout_id = ? AND user_id = ?", workoutID, claims.ID).Delete(&models.Exercise{}).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	// Then delete the workout
	if err := db.Where("id = ? AND user_id = ?", workoutID, claims.ID).Delete(&models.Workout{}).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, "deleted")
}
