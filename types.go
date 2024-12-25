package main

import (
	"golang.org/x/crypto/bcrypt"
	"log"
	"math/rand"
	"strconv"
	"time"
)

type CreateAccountRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password  string `json:"password"`
	Email     string `json:"email"`
}
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type TransferRequest struct {
	ToAccountNumber string `json:"toAccountNumber"`
	Amount          int64  `json:"amount"`
}

type Account struct {
	ID            int64     `json:"id"`
	FirstName     string    `json:"firstName"`
	LastName      string    `json:"lastName"`
	Email         string    `json:"email"`
	Password      string    `json:"password"`
	Balance       int64     `json:"balance"`
	AccountNumber string    `json:"accountNumber"`
	CreatedAt     time.Time `json:"createdAt"`
}

func NewAccount(firstName, lastName, email, password string) (*Account, error) {
	hashedPass, err := HashPassword(password)
	if err != nil {
		log.Fatal("Error hashing the password")
	}
	return &Account{
		FirstName:     firstName,
		LastName:      lastName,
		Email:         email,
		Password:      hashedPass,
		Balance:       0,
		AccountNumber: strconv.Itoa(rand.Intn(10000)),
		CreatedAt:     time.Now(),
	}, nil
}
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
