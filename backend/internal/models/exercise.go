package models

import "time"

type Exercise struct {
	ID uint `gorm:"primaryKey"`

	UserID uint  `gorm:"not null"`
	User   *User `gorm:"foreignKey:UserID"`

	WorkoutID uint    `gorm:"not null"`
	Workout   Workout `gorm:"foreignKey:WorkoutID"`

	Name string `gorm:"not null"`
	Note string

	Order int `gorm:"not null"`

	Sets int `gorm:"not null"`
	Reps int `gorm:"not null"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
	Deleted   bool `gorm:"default:false"`
}

type ExerciseResponse struct {
	ID        uint `json:"id"`
	WorkoutID uint `json:"workout_id"`

	Name string `json:"name"`
	Note string `json:"note"`

	Order int `json:"order"`

	Sets int `json:"sets"`
	Reps int `json:"reps"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (e *Exercise) CreateResponse() ExerciseResponse {
	return ExerciseResponse{
		ID:        e.ID,
		WorkoutID: e.WorkoutID,

		Name: e.Name,
		Note: e.Note,

		Order: e.Order,

		Sets: e.Sets,
		Reps: e.Reps,

		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}

}

type ExerciseSession struct {
	ID uint `gorm:"primaryKey"`

	ExerciseID uint      `gorm:"not null"`
	Exercise   *Exercise `gorm:"foreignKey:ExerciseID"`

	WorkoutSessionID uint            `gorm:"not null"`
	WorkoutSession   *WorkoutSession `gorm:"foreignKey:WorkoutSessionID"`

	Completed bool
	Skiped    bool
	Active    bool

	SetsDone int
}

type ExerciseSessionResponse struct {
	ID uint `json:"id"`

	Exercise         ExerciseResponse `json:"exercise"`
	WorkoutSessionID uint             `json:"workout_session_id"`

	Completed bool `json:"completed"`
	Skiped    bool `json:"skiped"`
	Active    bool `json:"active"`

	SetsDone int `json:"sets_done"`
}

func (es *ExerciseSession) CreateResponse() ExerciseSessionResponse {
	return ExerciseSessionResponse{
		ID: es.ID,

		Exercise:         es.Exercise.CreateResponse(),
		WorkoutSessionID: es.WorkoutSessionID,

		Completed: es.Completed,
		Skiped:    es.Skiped,
		Active:    es.Active,

		SetsDone: es.SetsDone,
	}

}
