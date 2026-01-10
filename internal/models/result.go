package models

import "time"

type Result struct {
	ID        int                    `json:"id" db:"id"`
	JobID     int                    `json:"job_id" db:"job_id"`
	Data      map[string]interface{} `json:"data" db:"data"`
	CreatedAt time.Time              `json:"created_at" db:"created_at"`
}
