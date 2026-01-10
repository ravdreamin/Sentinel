package models
import ("time")

type User struct {
	ID int `json:"id" db:"id"`
	Email string `json:"email" db:"email"`
	PasswordHash string `json:"password_hash" db:"password_hash"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}