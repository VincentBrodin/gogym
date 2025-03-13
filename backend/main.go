package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"

	"backend/config"
	"backend/internal/api"
	"backend/internal/models"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load env

	// Init
	e := echo.New()
	db, err := gorm.Open(postgres.Open(config.ENV["DB_CONNECTION_STRING"]), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Register models
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Workout{})
	db.AutoMigrate(&models.Exercise{})

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
	group := e.Group("/api")

	restricted := group.Group("/restricted")
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(models.JwtUserClaims)
		},
		SigningKey: []byte(config.ENV["JWT_SECRET"]),
	}
	restricted.Use(echojwt.WithConfig(config))

	group.POST("/login", api.Login)
	group.POST("/register", api.Register)

	restricted.PATCH("/account", api.EditAccount)

	restricted.GET("/workout/:id", api.GetWorkout)
	restricted.PUT("/workout", api.AddWorkout)
	restricted.PATCH("/workout/:id", api.EditWorkout)
	restricted.DELETE("/workout/:id", api.DeleteWorkout)

	restricted.GET("/exercise/:id", api.GetExercise)
	restricted.PUT("/exercise", api.AddExercise)
	restricted.PATCH("/exercise/:id", api.EditExercise)
	restricted.DELETE("/exercise/:id", api.DeleteExercise)

	restricted.GET("/test", func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*models.JwtUserClaims)
		name := claims.Username
		return c.String(http.StatusOK, "Welcome "+name+"!")
	})

	// Start
	e.Logger.Fatal(e.Start(":8080"))
}
