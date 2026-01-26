package main

import (
	"fmt"
	"log"
	"sentinel/internal/database"
	"sentinel/internal/server"
	"sentinel/internal/worker"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

}
func main() {
	dbPool, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}

	defer dbPool.Close()
	fmt.Println("🚀 Sentinel Database Connection Established")

	workerPool := worker.New(dbPool, 3)
	workerPool.Run()

	fmt.Println("⚡ Creating jobs in DB and sending to workers...")

	r := gin.Default()

	r.POST("/upload", server.UploadHandler)

	fmt.Println("Listening at 8081:")
	if err := r.Run(":8081"); err != nil {
		log.Fatal("Server failed to start: ", err)
	}

}
