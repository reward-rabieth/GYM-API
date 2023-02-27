package types

import (
	"time"
)

type Gymmember struct {
	Name            string    `json:"name"`
	Age             int       `json:"age"`
	Gender          string    `json:"gender"`
	Height          float64   `json:"height"`
	Weight          float64   `json:"weight"`
	Membership      string    `json:"membership"`
	StartDate       time.Time `json:"start_date"`
	EndDate         time.Time `json:"end_date"`
	PersonalTrainer string    `json:"personal_trainer"`
}

type Exercise struct {
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	MuscleGroups    []string `json:"muscle_groups"`
	EquipmnetNeeded []string `json:"equipment_needed"`
}

type Dietinformation struct {
	MealName string  `json:"meal_name"`
	Calories float64 `json:"calories"`
	Protein  float64 `json:"protein"`
	Carbs    float64 `json:"carbs"`
	Fat      float64 `json:"fat"`
	Fiber    float64 `json:"fiber"`
}

type ClassInformation struct {
	Name        string `json:"name"`
	Instructor  string `json:"instructor"`
	Description string `json:"description"`
	StartTime   string `json:"start_time"`
	EndTime     string `json:"end_time"`
}
type BookingInformation struct {
	ID           string    `json:"id"`
	CustomerName string    `json:"customer_name"`
	Date         time.Time `json:"date"`
	TimeSlot     string    `json:"time_slot"`
	Service      string    `json:"service"`
	TrainerID    string    `json:"trainer_id"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CreateGymMemberRequest struct {
	Name       string  `json:"name"`
	Age        int     `json:"age"`
	Gender     string  `json:"gender"`
	Height     float64 `json:"height"`
	Weight     float64 `json:"weight"`
	Membership string  `json:"membership"`
	//StartDate       time.Time `json:"start_date"`
	//EndDate         time.Time `json:"end_date"`
	PersonalTrainer string `json:"personal_trainer"`
}

func NewGymMember(req CreateGymMemberRequest) (*Gymmember, error) {

	return &Gymmember{
			Name:       req.Name,
			Age:        req.Age,
			Gender:     req.Gender,
			Height:     req.Height,
			Weight:     req.Weight,
			Membership: req.Membership,
			StartDate:  time.Now(),
			//EndDate:         req.EndDate,
			PersonalTrainer: req.PersonalTrainer,
		},
		nil
}
