package models

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtUserClaims struct {
	ID       uint   `json:"id"`
	Username string `json:"uname"`
	jwt.RegisteredClaims
}

type User struct {
	ID uint `gorm:"primaryKey"`

	Username string `gorm:"uniqueIndex"`
	Email    string `gorm:"uniqueIndex"`
	Password string

	Workouts         []Workout         `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
	Exercises        []Exercise        `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
	WorkoutSessions  []WorkoutSession  `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`
	ExerciseSessions []ExerciseSession `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserResponse struct {
	ID uint `json:"id"`

	Username string `json:"uname"`
	Email    string `json:"email"`

	Workouts []WorkoutResponse `json:"workouts"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
