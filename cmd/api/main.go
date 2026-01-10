package main

import (
	"fmt"
	"log"
	"sentinel/internal/database"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

}
func main() {
	pool, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}

	defer pool.Close()
	fmt.Println("🚀 Sentinel Database Connection Established!")

}
