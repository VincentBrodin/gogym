package models

import (
	"time"
)

type Workout struct {
	ID uint `gorm:"primaryKey"`

	UserID uint  `gorm:"not null"`
	User   *User `gorm:"foreignKey:UserID"`

	Exercises       []Exercise       `gorm:"foreignKey:WorkoutID;constraint:OnDelete:CASCADE;"`
	WorkoutSessions []WorkoutSession `gorm:"foreignKey:WorkoutID;constraint:OnDelete:CASCADE;"`

	Name string `gorm:"not null"`
	Note string

	LastDone  time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time 
	Deleted   bool `gorm:"default:false"`
}

type WorkoutResponse struct {
	ID uint `json:"id"`

	Name string `json:"name"`
	Note string `json:"note"`

	LastDone  time.Time `json:"last_done"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Exercises []ExerciseResponse `json:"exercises"`
}

func (w *Workout) CreateResponse() WorkoutResponse {
	response := WorkoutResponse{
		ID: w.ID,

		Name: w.Name,
		Note: w.Note,

		LastDone:  w.LastDone,
		CreatedAt: w.CreatedAt,
		UpdatedAt: w.UpdatedAt,
		Exercises: []ExerciseResponse{},
	}

	if w.Exercises != nil {
		response.Exercises = make([]ExerciseResponse, len(w.Exercises))
		for i, exercise := range w.Exercises {
			response.Exercises[i] = exercise.CreateResponse()
		}
	}

	return response
}

type WorkoutSession struct {
	ID uint `gorm:"primaryKey"`

	UserID uint  `gorm:"not null"`
	User   *User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`

	WorkoutID uint     `gorm:"not null"`
	Workout   *Workout `gorm:"foreignKey:WorkoutID;constraint:OnDelete:CASCADE;"`

	ExerciseSessions []ExerciseSession `gorm:"foreignKey:WorkoutSessionID;constraint:OnDelete:CASCADE;"`
	Active           bool

	StartedAt time.Time
	EndedAt   time.Time
}

type WorkoutSessionResponse struct {
	ID uint `json:"id"`

	Workout WorkoutResponse `json:"workout"`

	ExerciseSessions []ExerciseSessionResponse `json:"exercise_sessions"`

	Active bool `json:"active"`

	StartedAt time.Time `json:"started_at"`
	EndedAt   time.Time `json:"endend_at"`
}

func (ws *WorkoutSession) CreateResponse() WorkoutSessionResponse {
	response := WorkoutSessionResponse{
		ID: ws.ID,

		Workout: ws.Workout.CreateResponse(),

		Active: ws.Active,

		StartedAt:        ws.StartedAt,
		EndedAt:          ws.EndedAt,
		ExerciseSessions: []ExerciseSessionResponse{},
	}

	if ws.ExerciseSessions != nil {
		response.ExerciseSessions = make([]ExerciseSessionResponse, len(ws.ExerciseSessions))
		for i, exerciseSession := range ws.ExerciseSessions {
			response.ExerciseSessions[i] = exerciseSession.CreateResponse()
		}
	}

	return response
}
