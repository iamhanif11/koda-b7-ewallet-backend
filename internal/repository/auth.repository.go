package repository

import (
	"context"


	"github.com/iamhanif11/ewallet-backend/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (ar *AuthRepository) AddUser(ctx context.Context, email, hashPassword string) (model.User, error) {
	sql := `WITh register AS(INSERT INTO users(email, password) VALUES ($1, $2) RETURNING id, email, created_at), create_wallet AS ( INSERT INTO wallet (user_id) SELECT id FROM register) SELECT id, email, created_at FROM register;`

	args := []any{email, hashPassword}

	var user model.User
	if err := ar.db.QueryRow(ctx, sql, args...).Scan(&user.Id, &user.Email, &user.CreatedAt); err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (ar *AuthRepository) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	sql := `
		SELECT email, password
		FROM users
		WHERE email = $1
	`
	args := []any{email}

	var user model.User
	if err := ar.db.QueryRow(ctx, sql, args...).Scan(&user.Email, &user.Password); err != nil {
		return model.User{}, err
	}
	return user, nil
}


