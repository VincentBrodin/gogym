package models

import "time"

type Exercise struct {
	ID uint `gorm:"primaryKey"`

	UserID uint  `gorm:"not null"`
	User   *User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`

	WorkoutID uint    `gorm:"not null"`
	Workout   Workout `gorm:"foreignKey:WorkoutID;constraint:OnDelete:CASCADE;"`

	Name string `gorm:"not null"`
	Note string

	Order int `gorm:"not null"`

	Sets          int `gorm:"not null;default:3"`
	Reps          int `gorm:"not null;default:8"`
	RepsInReserve int `gorm:"not null;default:0"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
	Deleted   bool `gorm:"default:false"`

	ExerciseSessions []ExerciseSession `gorm:"foreignKey:ExerciseID;constraint:OnDelete:CASCADE;"`
	ExerciseWeights  []ExerciseWeight  `gorm:"foreignKey:ExerciseID;constraint:OnDelete:CASCADE;"`
}

type ExerciseResponse struct {
	ID        uint `json:"id"`
	WorkoutID uint `json:"workout_id"`

	Name string `json:"name"`
	Note string `json:"note"`

	Order int `json:"order"`

	Sets          int `json:"sets"`
	Reps          int `json:"reps"`
	RepsInReserve int `json:"rir"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	ExerciseWeights []ExerciseWeight `gorm:"foreignKey:ExerciseID;constraint:OnDelete:CASCADE;"`
}

func (e *Exercise) CreateResponse() ExerciseResponse {
	return ExerciseResponse{
		ID:        e.ID,
		WorkoutID: e.WorkoutID,

		Name: e.Name,
		Note: e.Note,

		Order: e.Order,

		Sets:          e.Sets,
		Reps:          e.Reps,
		RepsInReserve: e.RepsInReserve,

		CreatedAt: e.CreatedAt,
		UpdatedAt: e.UpdatedAt,
	}

}

type ExerciseSession struct {
	ID uint `gorm:"primaryKey"`

	UserID uint  `gorm:"not null"`
	User   *User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`

	ExerciseID uint      `gorm:"not null"`
	Exercise   *Exercise `gorm:"foreignKey:ExerciseID;constraint:OnDelete:CASCADE;"`

	WorkoutSessionID uint            `gorm:"not null"`
	WorkoutSession   *WorkoutSession `gorm:"foreignKey:WorkoutSessionID;constraint:OnDelete:CASCADE;"`

	Completed bool
	Skiped    bool
	Active    bool

	SetsDone int

	ExerciseWeights []ExerciseWeight `gorm:"foreignKey:ExerciseSessionID;constraint:OnDelete:CASCADE;"`
}

type ExerciseSessionResponse struct {
	ID uint `json:"id"`

	Exercise         ExerciseResponse `json:"exercise"`
	WorkoutSessionID uint             `json:"workout_session_id"`

	Completed bool `json:"completed"`
	Skiped    bool `json:"skiped"`
	Active    bool `json:"active"`

	SetsDone        int                      `json:"sets_done"`
	ExerciseWeights []ExerciseWeightResponse `json:"exercise_weights"`
}

func (es *ExerciseSession) CreateResponse() ExerciseSessionResponse {

	response := ExerciseSessionResponse{
		ID: es.ID,

		Exercise:         es.Exercise.CreateResponse(),
		WorkoutSessionID: es.WorkoutSessionID,

		Completed: es.Completed,
		Skiped:    es.Skiped,
		Active:    es.Active,

		SetsDone:        es.SetsDone,
		ExerciseWeights: []ExerciseWeightResponse{},
	}
	if es.ExerciseWeights != nil {
		response.ExerciseWeights = make([]ExerciseWeightResponse, len(es.ExerciseWeights))
		for i := range len(es.ExerciseWeights) {
			response.ExerciseWeights[i] = es.ExerciseWeights[i].CreateResponse()
		}
	}

	return response

}

type ExerciseWeight struct {
	ID uint `gorm:"primaryKey"`

	UserID uint  `gorm:"not null"`
	User   *User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;"`

	ExerciseID uint      `gorm:"not null"`
	Exercise   *Exercise `gorm:"foreignKey:ExerciseID;constraint:OnDelete:CASCADE;"`

	ExerciseSessionID uint             `gorm:"not null"`
	ExerciseSession   *ExerciseSession `gorm:"foreignKey:ExerciseSessionID;constraint:OnDelete:CASCADE;"`

	Set    int     `gorm:"not null"`
	Weight float64 `gorm:"not null"`
}

type ExerciseWeightResponse struct {
	ID uint `json:"id"`

	ExerciseId        uint `json:"exercise_id"`
	ExerciseSessionID uint `json:"exercise_session_id"`

	Set    int     `json:"set"`
	Weight float64 `json:"weight"`
}

func (ew *ExerciseWeight) CreateResponse() ExerciseWeightResponse {
	return ExerciseWeightResponse{
		ID: ew.ID,

		ExerciseId:        ew.ExerciseID,
		ExerciseSessionID: ew.ExerciseSessionID,

		Set:    ew.Set,
		Weight: ew.Weight,
	}
}
