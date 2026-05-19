package repository

import (
	"context"

	"github.com/iamhanif11/ewallet-backend/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (ur *UserRepository) GetProfileById(ctx context.Context, userId int) (model.User, error) {
	sql := `
		SELECT fullname, email, picture
		FROM users
		WHERE id = $1;
	`
	args := []any{userId}

	var user model.User
	if err := ur.db.QueryRow(ctx, sql, args...).Scan(&user.Fullname, &user.Email, &user.Picture); err != nil {

		return model.User{}, err
	}
	return user, nil
}

func (ur *UserRepository) GetPinById(ctx context.Context, userId int) (model.User, error) {
	sql := `
		SELECT pin 
		FROM users
		WHERE id= $1;
	`

	args := []any{userId}

	var user model.User
	if err := ur.db.QueryRow(ctx, sql, args...).Scan(&user.Pin); err != nil {
		return model.User{}, err
	}
	return user, nil
}
