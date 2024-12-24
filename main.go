package main

import "log"

func main() {

	store, _ := NewPostgresStore()
	err := store.Init()
	if err != nil {
		log.Fatal(err)
	}
	server := NewAPIServer(":3000", store)
	server.Run()
}
