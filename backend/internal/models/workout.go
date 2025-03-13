package models

import (
	"time"
)

type Workout struct {
	ID uint `gorm:"primaryKey"`

	UserID uint  `gorm:"not null"`
	User   *User `gorm:"foreignKey:UserID"`

	Exercises []Exercise

	Name string `gorm:"not null"`
	Note *string

	LastDone  time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

type WorkoutResponse struct {
	ID uint `json:"id"`

	Name string  `json:"name"`
	Note *string `json:"note"`

	LastDone  time.Time `json:"last_done"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Exercises []ExerciseResponse `json:"exercises"`
}
