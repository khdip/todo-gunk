package main

import (
	"log"
	"todo-gunk/todo/storage/postgres"
)

func main() {
	if err := postgres.Migrate(); err != nil {
		log.Fatal(err)
	}
}
