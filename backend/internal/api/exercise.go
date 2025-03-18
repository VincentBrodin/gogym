package api

import (
	"backend/internal/models"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type AddExerciseForm struct {
	ID        uint    `json:"id"`
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
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
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
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}
	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*models.JwtUserClaims)

	db := c.Get("db").(*gorm.DB)

	var exercises []models.Exercise
	if err := db.Where("workout_id = ? AND user_id = ?", form.WorkoutID, claims.ID).Find(&exercises).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	exercise := models.Exercise{
		UserID:    claims.ID,
		WorkoutID: form.WorkoutID,

		Name: form.Name,
		Note: form.Note,

		Order: len(exercises),

		Sets: form.Sets,
		Reps: form.Reps,

		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	if err := db.Create(&exercise).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Send back exercise
	return c.JSON(http.StatusOK, exercise.CreateResponse())
}

func EditExercise(c echo.Context) error {
	exerciseIDStr := c.Param("id")
	exerciseID, err := strconv.ParseUint(exerciseIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	var form AddExerciseForm
	if err := c.Bind(&form); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*models.JwtUserClaims)

	db := c.Get("db").(*gorm.DB)

	// Grab exercise
	var exercise models.Exercise
	if err := db.Where("id = ? AND user_id = ?", exerciseID, claims.ID).First(&exercise).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	// Update
	exercise.Name = form.Name
	exercise.Note = form.Note
	exercise.Sets = form.Sets
	exercise.Reps = form.Reps
	exercise.UpdatedAt = time.Now().UTC()

	if err := db.Save(&exercise).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, exercise.CreateResponse())
}

func EditAllExercises(c echo.Context) error {
	workoutIDStr := c.Param("id")
	workoutID, err := strconv.ParseUint(workoutIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	var forms []AddExerciseForm
	if err := c.Bind(&forms); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	if len(forms) == 0 {
		return c.JSON(http.StatusBadRequest, "Nothing to update")
	}

	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*models.JwtUserClaims)

	db := c.Get("db").(*gorm.DB)

	// Grab all exercises
	var exercises []models.Exercise
	if err := db.Where("workout_id = ? AND user_id = ?", workoutID, claims.ID).Find(&exercises).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Create a map to sync incoming with db
	exerciseMap := make(map[uint]*models.Exercise)
	for i := range exercises {
		exerciseMap[exercises[i].ID] = &exercises[i]
	}

	for _, form := range forms {
		if exercise, ok := exerciseMap[form.ID]; ok {
			exercise.Name = form.Name
			exercise.Note = form.Note
			exercise.Order = form.Order
			exercise.Sets = form.Sets
			exercise.Reps = form.Reps
			exercise.UpdatedAt = time.Now().UTC()
		}
	}

	// Make sure the order is good
	sort.Slice(exercises, func(i, j int) bool {
		return exercises[i].Order < exercises[j].Order
	})

	// Update order and save
	err = db.Transaction(func(tx *gorm.DB) error {
		for i := range exercises {
			exercises[i].Order = i
			if err := tx.Save(&exercises[i]).Error; err != nil {
				return err // transaction will be rolled back
			}
		}
		return nil // commit the transaction
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Create a response back
	responses := make([]models.ExerciseResponse, len(exercises))
	for i := range responses {
		responses[i] = exercises[i].CreateResponse()
	}

	return c.JSON(http.StatusOK, responses)
}

func DeleteExercise(c echo.Context) error {
	exerciseIDStr := c.Param("id")
	exerciseID, err := strconv.ParseUint(exerciseIDStr, 10, 32)

	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	userToken := c.Get("user").(*jwt.Token)
	claims := userToken.Claims.(*models.JwtUserClaims)

	db := c.Get("db").(*gorm.DB)

	// Grab the exercise to remove
	var exercise models.Exercise
	if err := db.Where("id = ? AND user_id = ?", exerciseID, claims.ID).First(&exercise).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	if err := db.Delete(&exercise).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Grab all remaning exercises
	var exercises []models.Exercise
	if err := db.Where("workout_id = ? AND user_id = ?", exercise.WorkoutID, claims.ID).Find(&exercises).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Sort by order
	sort.Slice(exercises, func(i, j int) bool {
		return exercises[i].Order < exercises[j].Order
	})

	// Update order and save
	err = db.Transaction(func(tx *gorm.DB) error {
		for i := range exercises {
			exercises[i].Order = i
			if err := tx.Save(&exercises[i]).Error; err != nil {
				return err // transaction will be rolled back
			}
		}
		return nil // commit the transaction
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Create a response back
	responses := make([]models.ExerciseResponse, len(exercises))
	for i := range responses {
		responses[i] = exercises[i].CreateResponse()
		fmt.Printf("%s %d\n", responses[i].Name, responses[i].Order)
	}

	return c.JSON(http.StatusOK, responses)
}
