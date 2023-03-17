package storage

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
	"github.com/reward-rabieth/gym/types"
	"github.com/reward-rabieth/gym/utils"

	"go.uber.org/zap"
)

type Storage interface {
	CreateMember(*types.Gymmember) error
	GetMembers() ([]*types.Gymmember, error)
	GetExercises() ([]*types.Exercise, error)
	CreateExercise(*types.Exercise) error
	GetMemberByid(int) (*types.Gymmember, error)
}

type PostgresStorage struct {
	db *sql.DB
}

func NewPostgresStorage(cfg utils.DbConfig) (*PostgresStorage, error) {
	constr := fmt.Sprintf("user=%s dbname=%s password=%s sslmode=%s", cfg.User, cfg.Name, cfg.Password, cfg.SSLMode)
	utils.Logger.Info("Connection string created", zap.String("connection_string", constr))

	db, err := sql.Open("postgres", constr)

	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresStorage{
		db: db,
	}, nil

}

func (p *PostgresStorage) Init() error {
	if err := p.CreateMemberTable(); err != nil {
		return err
	}
	if err := p.CreateExerciseTable(); err != nil {
		return err
	}
	return nil

}

func (p *PostgresStorage) CreateMemberTable() error {

	query := `CREATE TABLE IF NOT EXISTS members(
		id serial primary key,
		number serial,
		name varchar(255),
		age int,
		gender varchar(255),
		height float,
		weight float,
		membership varchar(255),
		start_date timestamp,
		end_date timestamp,
		personal_trainer varchar(50)

		
		)`
	_, err := p.db.Query(query)
	return err

}
func (p *PostgresStorage) CreateMember(member *types.Gymmember) error {

	query := `INSERT INTO members 
	 (number, name,age,gender,height,weight,membership,start_date,end_date,personal_trainer)
	values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
	`
	_, err := p.db.Query(
		query,
		member.Number,
		member.Name,
		member.Age,
		member.Gender,
		member.Height,
		member.Weight,
		member.Membership,
		member.StartDate,
		member.EndDate,
		member.PersonalTrainer,
	)
	if err != nil {
		return err
	}
	return nil
}
func (P *PostgresStorage) GetMemberByid(id int) (*types.Gymmember, error) {

	rows, err := P.db.Query("select * from members where id=$1", id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return scanIntoMembers(rows)
	}
	return nil, fmt.Errorf("member %d not found", id)

}
func (p *PostgresStorage) GetMembers() ([]*types.Gymmember, error) {
	rows, err := p.db.Query("select * from members")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	members := []*types.Gymmember{}
	for rows.Next() {
		member, err := scanIntoMembers(rows)
		if err != nil {
			return nil, err
		}
		members = append(members, member)
	}
	return members, err
}

func (p *PostgresStorage) CreateExerciseTable() error {
	query := `CREATE TABLE IF NOT EXISTS exercises(
		name varchar(255),
		description varchar(255),
		muscle_groups varchar(255),
		equipment varchar(255),
		sets int,
		reps int
	);`
	_, err := p.db.Query(query)

	return err
}

func (p *PostgresStorage) CreateExercise(exercise *types.Exercise) error {
	query := `INSERT INTO exercises
	(name,description,muscle_groups,equipment,sets,reps)
    values($1,$2,$3,$4,$5,$6)
	`
	_, err := p.db.Query(query,
		exercise.Name,
		exercise.Description,
		strings.Join(exercise.MuscleGroups, ", "),
		strings.Join(exercise.Equipment, ", "),
		exercise.Sets,
		exercise.Reps,
	)
	if err != nil {
		return err
	}
	return nil
}
func (p *PostgresStorage) GetExercises() ([]*types.Exercise, error) {
	rows, err := p.db.Query("select * from exercises")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	exercis := []*types.Exercise{}

	for rows.Next() {
		exercise, err := scanIntoExercises(rows)
		if err != nil {
			return nil, err
		}
		exercis = append(exercis, exercise)
	}
	return exercis, err

}

type exerciseDB struct {
	ID           int
	Name         string
	Description  string
	MuscleGroups string
	Equipment    string
	Sets         int
	Reps         int
}

func scanIntoExercises(rows *sql.Rows) (*types.Exercise, error) {
	exerciseDB := new(exerciseDB)

	err := rows.Scan(
		&exerciseDB.ID,
		&exerciseDB.Name,
		&exerciseDB.Description,
		&exerciseDB.MuscleGroups,
		&exerciseDB.Equipment,
		&exerciseDB.Sets,
		&exerciseDB.Reps,
	)
	if err != nil {
		return nil, err
	}

	equipment := strings.Split(exerciseDB.Equipment, ", ")

	return &types.Exercise{
		//	ID:           exerciseDB.ID,
		Name:         exerciseDB.Name,
		Description:  exerciseDB.Description,
		MuscleGroups: strings.Split(exerciseDB.MuscleGroups, ", "),
		Equipment:    equipment,
		Sets:         exerciseDB.Sets,
		Reps:         exerciseDB.Reps,
	}, nil
}

func scanIntoMembers(rows *sql.Rows) (*types.Gymmember, error) {
	member := new(types.Gymmember)

	err := rows.Scan(
		&member.ID,
		&member.Number,
		&member.Name,
		&member.Age,
		&member.Gender,
		&member.Height,
		&member.Weight,
		&member.Membership,
		&member.StartDate,
		&member.EndDate,
		&member.PersonalTrainer,
	)
	return member, err

}
