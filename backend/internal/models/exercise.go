package models

import "time"

type Exercise struct {
	ID uint `gorm:"primaryKey"`

	UserID uint  `gorm:"not null"`
	User   *User `gorm:"foreignKey:UserID"`

	WorkoutID uint    `gorm:"not null"`
	Workout   Workout `gorm:"foreignKey:WorkoutID"`

	Name string `gorm:"not null"`
	Note *string

	Sets int `gorm:"not null"`
	Reps int `gorm:"not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

type ExerciseResponse struct {
	ID        uint `json:"id"`
	WorkoutID uint `json:"workout_id"`

	Name string  `json:"name"`
	Note *string `json:"note"`

	Sets int `json:"sets"`
	Reps int `json:"reps"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
