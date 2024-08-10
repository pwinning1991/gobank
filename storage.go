package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(*Account) error
	UpdateAccount(*Account) error
	GetAccountByID(int) (*Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=postgres dbname=postgres password=gobank sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresStore{db: db}, nil
}

func (s *PostgresStore) Init() error {
	return s.createAccountTable()
}

func (s *PostgresStore) createAccountTable() error {
	query := `create table if not exists account(
		id serial PRIMARY KEY,
		firstname varchar(50),
		lastname varchar(50),
		nunber serial,
		balance serial,
		created_at timestamp
)`
	_, err := s.db.Exec(query)
	return err

}

func (s *PostgresStore) CreateAccount(account *Account) error {
	sqlStatement := `insert into account
	(firstName, lastName, number, balance, created_at)
	values ($1, $2, $3, $4, $5)`

	resp, err := s.db.Exec(sqlStatement,
		account.FirstName,
		account.LastName,
		account.Number,
		account.Balance,
		account.CreatedAt)
	if err != nil {
		return err
	}
	fmt.Printf("CreateAccount %v\n", resp)

	return nil
}

func (s *PostgresStore) DeleteAccount(account *Account) error {
	return nil
}

func (s *PostgresStore) UpdateAccount(account *Account) error {
	return nil
}

func (s *PostgresStore) GetAccountByID(id int) (*Account, error) {
	return nil, nil
}
