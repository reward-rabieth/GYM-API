package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/reward-rabieth/gym/types"
	"github.com/reward-rabieth/gym/utils"

	"go.uber.org/zap"
)

type Storage interface {
	CreateMember(*types.Gymmember) error
	GetMembers() ([]*types.Gymmember, error)
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
	return p.CreateMemberTable()

}

func (p *PostgresStorage) CreateMemberTable() error {

	query := `CREATE TABLE IF NOT EXISTS members(
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
	(name,age,gender,height,weight,membership,start_date,end_date,personal_trainer)
	values($1,$2,$3,$4,$5,$6,$7,$8,$9)
	`
	_, err := p.db.Query(
		query,
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

func (p *PostgresStorage) GetMembers() ([]*types.Gymmember, error) {
	rows, err := p.db.Query("select from *members")
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

func scanIntoMembers(rows *sql.Rows) (*types.Gymmember, error) {
	member := new(types.Gymmember)

	err := rows.Scan(
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
