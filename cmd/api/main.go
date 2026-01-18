package main

import (
	"context"
	"fmt"
	"log"
	"sentinel/internal/database"
	"sentinel/internal/models"
	"sentinel/internal/worker"

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

	for i := 0; i < 10; i++ {

		url := "https://google.com"
		var jobId int
		query := "INSERT INTO jobs (url, status) VALUES ($1, 'pending') RETURNING id"

		err := dbPool.QueryRow(context.Background(), query, url).Scan(&jobId)
		if err != nil {
			log.Printf("Failed to insert job: %v\n", err)
			continue
		}

		job := models.Job{
			ID:  jobId,
			URL: url,
		}

		workerPool.JobChan <- job
	}
	close(workerPool.JobChan)
	workerPool.Wg.Wait()
	fmt.Println("🏁 All jobs finished!")

}
