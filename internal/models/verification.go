package models

import "time"

type Verification struct {
	UserId   int       `json:"user_id" db:"user_id"`
	Code     string    `json:"code" db:"code"`
	ExpireAt time.Time `json:"expire_at" db:"expire_at"`
}
