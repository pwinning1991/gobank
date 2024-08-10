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
	GetAccounts() ([]*Account, error)
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
		number serial,
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
	sqlStatement := `select * from account where id = $1`
	rows, err := s.db.Query(sqlStatement, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		return scanIntoAccount(rows)
	}
	return nil, fmt.Errorf("Account %d  not found", id)
}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	sqlStatement := `select id, firstName, lastName, number, balance, created_at from account`
	rows, err := s.db.Query(sqlStatement)
	if err != nil {
		return nil, err
	}
	accounts := make([]*Account, 0)
	for rows.Next() {
		account, err := scanIntoAccount(rows)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)

	}
	return accounts, nil
}

func scanIntoAccount(rows *sql.Rows) (*Account, error) {
	account := &Account{}
	err := rows.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Number,
		&account.Balance,
		&account.CreatedAt)
	return account, err
}
