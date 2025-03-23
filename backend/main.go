package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"

	"backend/internal/api"
	"backend/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	ConnectionString string
	JwtSecret        string
	ContentDir       string
}

func spawnServer(config Config) *echo.Echo {
	// Init
	e := echo.New()
	db, err := gorm.Open(postgres.Open(config.ConnectionString), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Register models
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Workout{})
	db.AutoMigrate(&models.Exercise{})
	db.AutoMigrate(&models.WorkoutSession{})
	db.AutoMigrate(&models.ExerciseSession{})

	// Middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOriginFunc: func(origin string) (bool, error) {
			u, err := url.Parse(origin)
			if err != nil {
				return false, err
			}
			fmt.Printf("Request from %s\n", origin)
			host := u.Hostname()
			return host == "localhost" || host == "127.0.0.1", nil
		}, AllowHeaders: []string{"*"},
		AllowMethods:     []string{"GET", "HEAD", "PUT", "PATCH", "POST", "DELETE", "OPTIONS"},
		AllowCredentials: true,
	}))

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("db", db)
			return next(c)
		}
	})

	// Routers
	e.Static("/", config.ContentDir)
	e.File("/", filepath.Join(config.ContentDir, "index.html"))

	group := e.Group("/api")

	restricted := group.Group("/restricted")
	jwtConfig := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(models.JwtUserClaims)
		},
		SigningKey: []byte(config.JwtSecret),
	}
	restricted.Use(echojwt.WithConfig(jwtConfig))

	group.POST("/login", api.Login)
	group.POST("/register", api.Register)

	restricted.GET("/token", api.GetToken)
	restricted.GET("/account", api.GetAccount)
	restricted.PATCH("/account", api.EditAccount)

	restricted.PUT("/session/:id", api.StartSession)
	restricted.GET("/session", api.GetCurrentSession)
	restricted.PATCH("/session/:id", api.EditSession)

	restricted.GET("/workouts", api.GetAllWorkouts)
	restricted.GET("/workout/:id", api.GetWorkout)
	restricted.PUT("/workout", api.AddWorkout)
	restricted.PATCH("/workout/:id", api.EditWorkout)
	restricted.DELETE("/workout/:id", api.DeleteWorkout)

	restricted.GET("/exercise/:id", api.GetExercise)
	restricted.PUT("/exercise", api.AddExercise)
	restricted.PATCH("/exercise/:id", api.EditExercise)
	restricted.PATCH("/exercises/:id", api.EditAllExercises)
	restricted.DELETE("/exercise/:id", api.DeleteExercise)

	restricted.GET("/test", func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*models.JwtUserClaims)
		name := claims.Username
		return c.String(http.StatusOK, "Welcome "+name+"!")
	})

	// Start
	return e
}

func main() {
	ex, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	appDir := filepath.Dir(ex)
	envPath := filepath.Join(appDir, ".env")

	if err := godotenv.Load(envPath); err != nil {
		log.Printf("Error loading .env from %s: %v", envPath, err)
	}

	serverConfig := Config{
		ConnectionString: os.Getenv("DB_CONNECTION_STRING"),
		JwtSecret:        os.Getenv("JWT_SECRET"),
		ContentDir:       filepath.Join(appDir, "content"),
	}

	fmt.Println(serverConfig.ConnectionString)

	e := spawnServer(serverConfig)
	e.Logger.Fatal(e.Start(":8080"))

}
