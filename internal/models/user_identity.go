package models

import "time"


type UserIdentity struct {
	ID         int       `json:"id" db:"id"`
	UserID     int       `json:"user_id" db:"user_id"`
	Provider   string    `json:"provider" db:"provider"`       // e.g., "local", "google"
	ProviderID string    `json:"provider_id" db:"provider_id"` // OAuth subject ID or empty for local
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}
