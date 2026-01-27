package database

import (
	"context"
	"fmt"
	"sentinel/internal/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateUser(p *pgxpool.Pool, u *models.User) error {
	query := "INSERT INTO users(email,password_hash) VALUES ($1,$2) RETURNING id, created_at, is_verified"
	err := p.QueryRow(context.Background(), query, u.Email, u.PasswordHash).Scan(&u.ID, u.CreatedAt, u.IsVerified)
	if err != nil {
		return fmt.Errorf("unable to insert user: %w", err)

	}
	return nil
}

func SaveVerification(p *pgxpool.Pool, v *models.Verification) error {
	query := `INSERT INTO verifications (user_id, code, expires_at)
		          VALUES ($1, $2, $3)
		          ON CONFLICT (user_id)
		          DO UPDATE SET code = EXCLUDED.code, expires_at = EXCLUDED.expires_at`
	_, err := p.Exec(context.Background(), query, v.UserId, v.Code, v.ExpireAt)
	if err != nil {
		return fmt.Errorf("Error storing verification: %s", err)

	}
	return nil
}

func GetUserByEmail(pool *pgxpool.Pool, email string) (*models.User, error) {

	var user models.User
	query := `SELECT id, email, password_hash, is_verified FROM users WHERE email = $1`
	err := pool.QueryRow(context.Background(), query, email).
		Scan(&user.ID, &user.Email, &user.PasswordHash, &user.IsVerified)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetVerification(pool *pgxpool.Pool, userID int) (*models.Verification, error) {
	var v models.Verification
	query := `SELECT user_id, code, expires_at FROM verifications WHERE user_id = $1`
	err := pool.QueryRow(context.Background(), query, userID).
		Scan(&v.UserId, &v.Code, &v.ExpireAt)

	if err != nil {
		return nil, err
	}
	return &v, nil
}

func MarkUserVerified(pool *pgxpool.Pool, userID int) error {
	query := `UPDATE users SET is_verified = TRUE WHERE id = $1`
	_, err := pool.Exec(context.Background(), query, userID)
	return err
}

func DeleteVerification(pool *pgxpool.Pool, userID int) error {
	query := `DELETE FROM verifications WHERE user_id = $1`
	_, err := pool.Exec(context.Background(), query, userID)
	return err
}
