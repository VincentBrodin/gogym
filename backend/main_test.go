package main

import (
	"backend/internal/api"
	"backend/internal/models"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	dummyJwt      = ""
	dummyUser     models.User
	dummyWorkout  models.Workout
	dummyExercise models.Exercise
)

func resetDatabase() error {
	db, err := gorm.Open(postgres.Open(os.Getenv("DB_CONNECTION_STRING_TEST")), &gorm.Config{})
	if err != nil {
		return err
	}

	if err := db.Migrator().DropTable(&models.User{}, &models.Workout{}, &models.Exercise{}); err != nil {
		return err
	}
	if err := db.AutoMigrate(&models.User{}, &models.Workout{}, &models.Exercise{}); err != nil {
		return err
	}

	bytes, err := bcrypt.GenerateFromPassword([]byte("dummypswd"), 14)
	if err != nil {
		return err
	}

	dummyUser = models.User{
		Username:  "dummy",
		Email:     "dummy@test.com",
		Password:  string(bytes),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	if err = db.Create(&dummyUser).Error; err != nil {
		return err
	}
	dummyJwt, err = api.GenerateUserToken(dummyUser)

	dummyWorkout = models.Workout{
		UserID:    dummyUser.ID,
		Name:      "dummy workout",
		Note:      "some info about dummy workout",
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		LastDone:  time.Now().UTC(),
	}
	if err = db.Create(&dummyWorkout).Error; err != nil {
		return err
	}

	dummyExercise = models.Exercise{
		UserID:    dummyUser.ID,
		WorkoutID: dummyWorkout.ID,
		Name:      "dummy exercise",
		Note:      "some info about dummy workout",
		Sets:      2,
		Reps:      8,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}
	if err = db.Create(&dummyExercise).Error; err != nil {
		return err
	}

	return err
}

func loadEnv() error {
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Join(filepath.Dir(filename), ".")
	envPath := filepath.Join(dir, ".env")
	if err := godotenv.Load(envPath); err != nil {
		return err
	}
	return nil
}

func sendTestForm(endpoint, method string, payload url.Values) (*httptest.ResponseRecorder, error) {
	if err := loadEnv(); err != nil {
		return nil, err
	}

	testConfig := Config{
		ConnectionString: os.Getenv("DB_CONNECTION_STRING_TEST"),
		JwtSecret:        os.Getenv("JWT_SECRET"),
	}

	e := spawnServer(testConfig)

	if err := resetDatabase(); err != nil {
		return nil, err
	}

	req := httptest.NewRequest(method, endpoint, strings.NewReader(payload.Encode()))
	req.Header.Set(echo.HeaderContentType, "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Bearer "+dummyJwt)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)
	return rec, nil
}

func sendTestJson(endpoint, method string, payload any) (*httptest.ResponseRecorder, error) {
	if err := loadEnv(); err != nil {
		return nil, err
	}
	testConfig := Config{
		ConnectionString: os.Getenv("DB_CONNECTION_STRING_TEST"),
		JwtSecret:        os.Getenv("JWT_SECRET"),
	}

	e := spawnServer(testConfig)

	if err := resetDatabase(); err != nil {
		return nil, err
	}

	reqBody, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req := httptest.NewRequest(method, endpoint, bytes.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set("Authorization", "Bearer "+dummyJwt)
	rec := httptest.NewRecorder()

	e.ServeHTTP(rec, req)
	return rec, nil
}

// Auth
func TestLoginEndpoint(t *testing.T) {
	form := url.Values{}
	form.Set("uname", "dummy2")
	form.Set("pswd", "dummypswd")

	rec, err := sendTestForm("/api/restricted/account", http.MethodPatch, form)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

}

func TestRegisterEndpoint(t *testing.T) {
	form := url.Values{}
	form.Set("uname", "dummy2")
	form.Set("email", "dummy2@test.com")
	form.Set("pswd", "dummypswd")

	rec, err := sendTestForm("/api/restricted/account", http.MethodPatch, form)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestEditAccount(t *testing.T) {
	form := url.Values{}
	form.Set("uname", "dummy edited")
	rec, err := sendTestForm("/api/restricted/account", http.MethodPatch, form)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

// Workout
func TestAddWorkout(t *testing.T) {
	payload := api.AddWorkoutForm{
		Name: "dummy workout",
	}

	rec, err := sendTestJson("/api/restricted/workout", http.MethodPut, payload)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestEditWorkout(t *testing.T) {
	payload := api.AddWorkoutForm{
		Name: "dummy workout edit",
	}

	rec, err := sendTestJson(fmt.Sprintf("/api/restricted/workout/%d", dummyWorkout.ID), http.MethodPatch, payload)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestRemoveWorkout(t *testing.T) {
	rec, err := sendTestJson(fmt.Sprintf("/api/restricted/workout/%d", dummyWorkout.ID), http.MethodDelete, nil)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

// Exercise
func TestAddExercise(t *testing.T) {
	payload := api.AddExerciseForm{
		WorkoutID: dummyWorkout.ID,

		Name: "dummy exercise",
		Note: "dummy note",
		Sets: 2,
		Reps: 8,
	}

	rec, err := sendTestJson("/api/restricted/exercise", http.MethodPut, payload)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestEditExercise(t *testing.T) {
	payload := api.AddExerciseForm{
		Name: "dummy exercise edit",
		Note: "edit note",
		Sets: 3,
		Reps: 12,
	}

	rec, err := sendTestJson(fmt.Sprintf("/api/restricted/exercise/%d", dummyExercise.ID), http.MethodPatch, payload)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}

func TestRemoveExercise(t *testing.T) {
	rec, err := sendTestJson(fmt.Sprintf("/api/restricted/exercise/%d", dummyExercise.ID), http.MethodDelete, nil)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
}
