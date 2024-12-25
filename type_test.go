package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewAccount(t *testing.T) {
	acc, err := NewAccount("Omar", "a", "omar@gg.com", "Aa123")
	assert.Nil(t, err)
	fmt.Printf("%+v\n", acc)

}
