package main

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"reflect"
	"strconv"
)

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}
func permissionDenied(w http.ResponseWriter) {
	WriteJSON(w, http.StatusUnauthorized, ApiErr{Error: "invalid token"})
}
func UseJWT(handlerFunc http.HandlerFunc, s Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("calling JWT auth middleware")

		tokenString := r.Header.Get("x-jwt-token")
		token, err := validateJWT(tokenString)
		if err != nil {
			fmt.Print("error in validateJWT")
			permissionDenied(w)
			return
		}
		if !token.Valid {
			fmt.Print("here2")
			permissionDenied(w)
			return
		}
		userID, err := getID(r)
		if err != nil {
			fmt.Print("here")
			permissionDenied(w)
			return
		}
		account, err := s.GetAccountById(userID)
		if err != nil {
			permissionDenied(w)
			fmt.Print("here")
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		cAcc := strconv.FormatFloat(claims["accountNumber"].(float64), 'f', -1, 64)
		println(cAcc, reflect.TypeOf(cAcc))
		println(account.AccountNumber, reflect.TypeOf(account.AccountNumber))
		if account.AccountNumber != cAcc {
			permissionDenied(w)
			return
		}

		if err != nil {
			WriteJSON(w, http.StatusForbidden, ApiErr{Error: "invalid token"})
			return
		}

		handlerFunc(w, r)
	}
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secret), nil
	})
}

func CreateJWTToken(acc *Account) (string, error) {
	claims := jwt.MapClaims{
		"expiresAt":     1500,
		"accountNumber": acc.AccountNumber,
	}
	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))

}

func getID(r *http.Request) (int, error) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return id, fmt.Errorf("invalid id given %s", idStr)
	}
	return id, nil
}
