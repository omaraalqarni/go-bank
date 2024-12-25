package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	UpdateAccount(*Account) (*Account, error)
	GetAccountById(int) (*Account, error)
	GetAllAccounts() ([]*Account, error)
	DeleteAccountById(int) error
	LoginUser(LoginRequest) (*Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

func (p *PostgresStore) Init() error {
	return p.CreateAccountTable()
}

func (p *PostgresStore) CreateAccountTable() error {
	query := `CREATE TABLE IF NOT EXISTS account(
		id serial primary key,
		first_name varchar(255),
		last_name varchar(255),
    	email varchar(255),
    	password varchar(255),
		balance bigint,
    	account_number varchar(255),
		created_at timestamp
    );`

	_, err := p.db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgresStore) DropTablesCascade() error {
	query := "DROP TABLE account;"
	if _, err := p.db.Exec(query); err != nil {
		return err
	}
	return nil
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "user=root dbname=postgres password=root sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return &PostgresStore{
		db: db,
	}, nil
}

func (p PostgresStore) CreateAccount(account *Account) error {
	query := `INSERT INTO account (first_name, last_name, email, password, balance, account_number, created_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7) `

	_, err := p.db.Query(query, account.FirstName, account.LastName, account.Email, account.Password, account.Balance, account.AccountNumber, account.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (p *PostgresStore) LoginUser(loginReq LoginRequest) (*Account, error) {
	row, err := p.db.Query("SELECT * WHERE email == $1;", loginReq.Email)
	if err != nil {
		return nil, err
	}
	account, err := scanIntoAccount(row)
	if err != nil {
		return nil, err
	}
	if !CheckPasswordHash(loginReq.Password, account.Password) {
		return nil, err
	}
	return nil, nil
}

func (p *PostgresStore) UpdateAccount(account *Account) (*Account, error) {
	//TODO implement me
	return nil, nil
}

func (p *PostgresStore) GetAccountById(id int) (*Account, error) {

	res, err := p.db.Query("select * from account where id =$1", id)
	if err != nil {
		return nil, err
	}
	for res.Next() {
		return scanIntoAccount(res)
	}
	return nil, fmt.Errorf("account %d not found", id)
}

func (p *PostgresStore) GetAllAccounts() ([]*Account, error) {
	query := ` select * from account;`
	rows, err := p.db.Query(query)
	if err != nil {
		return nil, err
	}
	accounts := []*Account{}
	defer rows.Close()
	for rows.Next() {
		account, err := scanIntoAccount(rows)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, account)
	}
	return accounts, nil
}
func (p *PostgresStore) DeleteAccountById(id int) error {
	if _, err := p.db.Query("delete from account where id=$1", id); err != nil {
		return err
	}
	return nil
}
func scanIntoAccount(rows *sql.Rows) (*Account, error) {
	account := new(Account)
	err := rows.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Email,
		&account.Balance,
		&account.AccountNumber,
		&account.CreatedAt,
	)
	return account, err
}
