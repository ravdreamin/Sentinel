package database

import (
	"context"
	"encoding/json"
	"path/filepath"
	"sentinel/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CountUserFiles(pool *pgxpool.Pool, userID int) (int, error) {
	query := "SELECT COUNT(DISTINCT file_path) FROM jobs WHERE user_id = $1"
	var count int
	err := pool.QueryRow(context.Background(), query, userID).Scan(&count)
	return count, err
}

func GetJobProgress(pool *pgxpool.Pool, filePath string) (total int, completed int, failed int, err error) {
	queryTotal := "SELECT COUNT(*) FROM jobs WHERE file_path = $1"
	queryCompleted := "SELECT COUNT(*) FROM jobs WHERE file_path = $1 AND status = 'Completed'"
	// We don't have 'Failed' status yet in worker, but we will add it.
	queryFailed := "SELECT COUNT(*) FROM jobs WHERE file_path = $1 AND status = 'Failed'"

	if err := pool.QueryRow(context.Background(), queryTotal, filePath).Scan(&total); err != nil {
		return 0, 0, 0, err
	}
	if err := pool.QueryRow(context.Background(), queryCompleted, filePath).Scan(&completed); err != nil {
		return 0, 0, 0, err
	}
	if err := pool.QueryRow(context.Background(), queryFailed, filePath).Scan(&failed); err != nil {
		return 0, 0, 0, err
	}
	return total, completed, failed, nil
}

func GetJobResults(pool *pgxpool.Pool, filePath string) ([]models.CrawlData, error) {
	// Join jobs and results
	query := `
        SELECT r.data 
        FROM results r 
        JOIN jobs j ON r.job_id = j.id 
        WHERE j.file_path = $1
    `
	rows, err := pool.Query(context.Background(), query, filePath)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.CrawlData
	for rows.Next() {
		var dataJSON []byte
		if err := rows.Scan(&dataJSON); err != nil {
			continue
		}
		var data models.CrawlData
		if err := json.Unmarshal(dataJSON, &data); err == nil {
			results = append(results, data)
		}
	}
	return results, nil
}

func DeleteJobByFilePath(pool *pgxpool.Pool, filePath string, userID int) error {
	// Delete jobs, cascading will delete results.
	// Ensure it belongs to user
	query := "DELETE FROM jobs WHERE file_path = $1 AND user_id = $2"
	_, err := pool.Exec(context.Background(), query, filePath, userID)
	return err
}

func GetUserJobs(pool *pgxpool.Pool, userID int) ([]string, error) {
	query := "SELECT DISTINCT file_path FROM jobs WHERE user_id = $1 ORDER BY file_path DESC"
	rows, err := pool.Query(context.Background(), query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []string
	for rows.Next() {
		var f string
		if err := rows.Scan(&f); err == nil {
			// Extract just the filename from the path (e.g., ./uploads/filename -> filename)
			files = append(files, filepath.Base(f))
		}
	}
	return files, nil
}
