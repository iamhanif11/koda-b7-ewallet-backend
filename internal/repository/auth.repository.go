package repository

import (
	"context"
	"log"
	"time"

	"github.com/iamhanif11/ewallet-backend/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type AuthRepository struct {
	db  *pgxpool.Pool
	rdb *redis.Client
}

func NewAuthRepository(db *pgxpool.Pool, rdb *redis.Client) *AuthRepository {
	return &AuthRepository{
		db:  db,
		rdb: rdb,
	}
}

func (ar *AuthRepository) AddUser(ctx context.Context, email, hashPassword string) (model.User, error) {
	sql := `WITh register AS(
	INSERT INTO users(email, password) VALUES ($1, $2) 
	RETURNING id, email, created_at
	), 
	create_wallet AS ( 
	INSERT INTO wallet (user_id) SELECT id FROM register
	) 
	SELECT id, email, created_at FROM register;`
	log.Println(email, hashPassword)
	args := []any{email, hashPassword}

	var user model.User
	if err := ar.db.QueryRow(ctx, sql, args...).Scan(&user.Id, &user.Email, &user.CreatedAt); err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (ar *AuthRepository) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	log.Println(email)
	sql := `
		SELECT id, email, password
		FROM users
		WHERE email = $1
	`
	args := []any{email}

	var user model.User
	if err := ar.db.QueryRow(ctx, sql, args...).Scan(&user.Id, &user.Email, &user.Password); err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (ar *AuthRepository) BlacklistToken(ctx context.Context, token string, expired time.Duration) error {
	return ar.rdb.Set(ctx, token, "revoked", expired).Err()
}

func (ar *AuthRepository) IsTokenBlacklisted(ctx context.Context, token string) bool {
	err := ar.rdb.Get(ctx, token).Err()

	return err == nil
}

func (ar *AuthRepository) CheckPinUserByEmail(ctx context.Context, email string) (bool, error) {
	sql := `
		SELECT pin
		FROM users 
		WHERE email = $1
	`

	var result string
	err := ar.db.QueryRow(ctx, sql, email).Scan(&result)
	if err != nil {
		return false, err
	}
	return true, nil
}
