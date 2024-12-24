package main

import (
	"math/rand"
	"strconv"
	"time"
)

type CreateAccountRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type TransferRequest struct {
	ToIBAN string `json:"to_iban"`
	Amount int64  `json:"amount"`
}

type Account struct {
	ID        int64     `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Balance   int64     `json:"balance"`
	IBAN      string    `json:"iban"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewAccount(firstName, lastName string) *Account {
	return &Account{
		FirstName: firstName,
		LastName:  lastName,
		Balance:   0,
		IBAN:      strconv.Itoa(rand.Intn(10000)),
		CreatedAt: time.Now(),
	}
}
