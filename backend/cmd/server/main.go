package main

import (
	"fmt"
	"log"

	"url-shortener-backend/internal/repository"
)

func main() {
	db, err := repository.InitDB()
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	defer db.Close()

	fmt.Println("Database connected successfully!")

	err = repository.CreateTable(db)
	if err != nil {
		log.Fatal("Table creation failed:", err)
	}

	fmt.Println("Links table created successfully!")
}
