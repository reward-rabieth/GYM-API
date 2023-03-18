package types

import (
	"math/rand"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Number   int64  `json:"number"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Number int64  `json:"number"`
	Token  string `json:"token"`
}

type Gymmember struct {
	ID                int       `json:"id"`
	Number            int64     `json:"membership_no"`
	Name              string    `json:"name"`
	EncryptedPassword string    `json:"-"`
	Age               int       `json:"age"`
	Gender            string    `json:"gender"`
	Height            float64   `json:"height"`
	Weight            float64   `json:"weight"`
	Membership        string    `json:"membership"`
	StartDate         time.Time `json:"start_date"`
	EndDate           time.Time `json:"end_date"`
	PersonalTrainer   string    `json:"personal_trainer"`
}

func (a *Gymmember) ValidatePassword(pwd string) bool {
	return bcrypt.CompareHashAndPassword([]byte(a.EncryptedPassword), []byte(pwd)) == nil

}

type Exercise struct {
	//ID           int      `json:"id"`
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	MuscleGroups []string `json:"muscle_groups"`
	Equipment    []string `json:"equipment"`
	Sets         int      `json:"sets"`
	Reps         int      `json:"reps"`
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
	Password   string  `json:"password"`
	//StartDate     time.Time `json:"start_date"`
	//EndDate         time.Time `json:"end_date"`
	PersonalTrainer string `json:"personal_trainer"`
}

type GymParams struct {
	Name       string
	Age        int
	Gender     string
	Height     float64
	Weight     float64
	Membership string

	//StartDate     time.Time `json:"start_date"`
	//EndDate         time.Time `json:"end_date"`
	PersonalTrainer string
}

func NewGymMember(req GymParams, password string) (*Gymmember, error) {
	encrpwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &Gymmember{
			Name:              req.Name,
			Age:               req.Age,
			Number:            int64(rand.Intn(1000)),
			EncryptedPassword: string(encrpwd),
			Gender:            req.Gender,
			Height:            req.Height,
			Weight:            req.Weight,
			Membership:        req.Membership,
			StartDate:         time.Now(),
			//EndDate:         req.EndDate,
			PersonalTrainer: req.PersonalTrainer,
		},
		nil
}

type CreateNewExerciseRequest struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	MuscleGroups []string `json:"muscle_groups"`
	Equipment    []string `json:"-"`
	EquipmentStr string   `json:"equipment"`
	Sets         int      `json:"sets"`
	Reps         int      `json:"reps"`
}

func NewExercise(req CreateNewExerciseRequest) (*Exercise, error) {
	equipment := strings.Split(req.EquipmentStr, ",")
	for i, e := range equipment {
		equipment[i] = strings.TrimSpace(e)
	}
	return &Exercise{
		Name:         req.Name,
		Description:  req.Description,
		MuscleGroups: req.MuscleGroups,
		Equipment:    equipment,
		Sets:         req.Sets,
		Reps:         req.Reps,
	}, nil
}
