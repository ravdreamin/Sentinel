package models

import "time"

type Job struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	URL       string    `json:"url" db:"url"`
	FilePath  string    `json:"file_path,omitempty" db:"file_path"`
	JobType   string    `json:"job_type" db:"job_type"`
	Status    string    `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
