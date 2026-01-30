package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Metric struct for aggregated data
type JobMetrics struct {
	TotalRequests       int            `json:"total_requests"`
	AverageResponseTime float64        `json:"avg_response_time"`
	StatusCodes         map[string]int `json:"status_codes"`
	TotalDataSize       int            `json:"total_data_size"` // Estimated
}

func GetJobMetrics(pool *pgxpool.Pool, filePath string) (JobMetrics, error) {
	metrics := JobMetrics{StatusCodes: make(map[string]int)}

	// We fetch all results and aggregate in Go to avoid complex SQL for now,
	// or use smart SQL. Let's use SQL for efficiency where possible but we have JSONB.
	// Casting JSONB to int in Postgres: (data->>'response_time')::int

	// Metrics: Total Requests, Avg Response Time, Total Data Size
	queryStats := `
        SELECT 
            COUNT(*), 
            COALESCE(AVG((r.data->>'response_time')::int), 0),
            COALESCE(SUM(OCTET_LENGTH(r.data::text)), 0)
        FROM results r 
        JOIN jobs j ON r.job_id = j.id 
        WHERE j.file_path = $1
    `
	err := pool.QueryRow(context.Background(), queryStats, filePath).Scan(&metrics.TotalRequests, &metrics.AverageResponseTime, &metrics.TotalDataSize)
	if err != nil {
		return metrics, err
	}

	// Status Codes Distribution
	queryCodes := `
        SELECT 
            r.data->>'status_code' as code, 
            COUNT(*) 
        FROM results r 
        JOIN jobs j ON r.job_id = j.id 
        WHERE j.file_path = $1 
        GROUP BY 1
    `
	rows, err := pool.Query(context.Background(), queryCodes, filePath)
	if err != nil {
		return metrics, err // Return partial metrics if code query fails
	}
	defer rows.Close()

	for rows.Next() {
		var code string
		var count int
		if err := rows.Scan(&code, &count); err == nil {
			metrics.StatusCodes[code] = count
		}
	}

	return metrics, nil
}
