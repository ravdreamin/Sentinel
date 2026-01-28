package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"sentinel/internal/models"
)

func CreateJob(dbPool *pgxpool.Pool, job *models.Job) error {
	query := "INSERT INTO jobs(url,status,file_path,job_type,user_id) VALUES ($1,$2,$3,$4,$5) RETURNING id"
	err := dbPool.QueryRow(context.Background(), query, job.URL, job.Status, job.FilePath, job.JobType, job.UserID).Scan(&job.ID)
	if err != nil {
		fmt.Printf("Failed to insert job: %v\n", err)

	}

	return err
}

func UpdateJobStatus(pool *pgxpool.Pool, jobId int, status string) error {
	query := "UPDATE jobs SET status = $1 WHERE id = $2"
	_, err := pool.Exec(context.Background(), query, status, jobId)
	if err != nil {
		fmt.Println("Error updating jobs", err)

	}

	return err

}
