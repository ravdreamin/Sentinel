package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	"sentinel/internal/models"
)

func CreateJob(dbPool *pgxpool.Pool, job models.Job) error {
	query := "INSERT INTO jobs(url,status,file_path,job_type) VALUES ($1,$2,$3,$4) RETURNING id"
	err := dbPool.Exec(context.Background(), query, job.URL, job.Status, job.FilePath, job.JobType).Scan(&job.ID)
	if err != nil {
		log.Printf("Failed to insert job: %v\n", err)
	}

	return err
}
