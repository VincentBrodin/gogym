package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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
	db.AutoMigrate(&models.User{}, &models.Workout{}, &models.WorkoutSession{}, &models.Exercise{}, &models.ExerciseSession{}, &models.ExerciseWeight{})

	// Middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOriginFunc: func(origin string) (bool, error) {
			// u, err := url.Parse(origin)
			// if err != nil {
			// 	return false, err
			// }
			// host := u.Hostname()
			// return host == "localhost" || host == "127.0.0.1", nil
			return true, nil
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
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root:  config.ContentDir,
		Index: "index.html",
		HTML5: true,
	}))

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
	restricted.DELETE("/account", api.DeleteAccount)

	restricted.PUT("/session/:id", api.StartSession)
	restricted.GET("/session", api.GetCurrentSession)
	restricted.GET("/session/:id", api.GetSession)
	restricted.GET("/sessions", api.GetAllSessions)
	restricted.PATCH("/session/:id", api.EditSession)
	restricted.DELETE("/session/:id", api.DeleteSession)

	restricted.GET("/workouts", api.GetAllWorkouts)
	restricted.GET("/workout/:id", api.GetWorkout)
	restricted.PUT("/workout", api.AddWorkout)
	restricted.PUT("/workout/copy/:id", api.CopyWorkout)
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

	e.GET("/*", func(c echo.Context) error {
		if strings.HasPrefix(c.Request().URL.Path, "/api") {
			return echo.ErrNotFound
		}
		return c.File(filepath.Join(config.ContentDir, "index.html"))
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
