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

func (ur *UserRepository) UpdateProfileById(ctx context.Context, userId int, fullname, phone, picture *string) (model.User, error) {
	sql := `
		UPDATE users
		SET
			fullname = $2,
			phone = $3,
			picture = $4
		WHERE id = $1
		RETURNING id, fullname, phone, picture;
	`

	args := []any{userId, fullname, phone, picture}

	var user model.User
	if err := ur.db.QueryRow(ctx, sql, args...).Scan(&user.Id, &user.Fullname, &user.Picture, &user.Phone); err != nil {
		return model.User{}, err
	}
	return user, nil
}
func (ur *UserRepository) GetPasswordById(ctx context.Context, userId int) (model.User, error) {
	sql := `
		SELECT password
		FROM users
		WHERE id = $1;
	`

	args := []any{userId}

	var user model.User
	if err := ur.db.QueryRow(ctx, sql, args...).Scan(&user.Password); err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (ur *UserRepository) UpdatePasswordById(ctx context.Context, userId int, hashedPassword string) error {
	sql := `
		UPDATE users
		SET
		password = $2,
		updated_at = NOW()
		WHERE id = $1;
	`
	args := []any{userId, hashedPassword}

	_, err := ur.db.Exec(ctx, sql, args...)
	return err
}

func (ur *UserRepository) UpdatedPinById(ctx context.Context, userId int, pin string) error {
	sql := `
		UPDATE users
		SET
			pin = $2, 
			updated_at = NOW()
		WHERE id = $1		
	`
	args := []any{userId, pin}

	_, err := ur.db.Exec(ctx, sql, args...)
	return err
}

func (ur *UserRepository) GetDashboardInformationById(ctx context.Context, userId int) (model.UserDashboardInformation, error) {
	sql := `
	SELECT
		w.balance AS balance,
		COALESCE(SUM(
			CASE
				WHEN t.status = 'completed' AND t.transaction_type IN ('top-up', 'transfer in')
				THEN t.amount
				ELSE 0
			END	
		), 0) AS income,
		COALESCE(SUM(
			CASE
				WHEN t.status = 'completed' AND t.transaction_type = 'transfer out'
				THEN t.amount
				ELSE 0
			END	
		), 0) AS expense
	FROM wallet w
	LEFT JOIN transactions t ON t.user_id = w.user_id
	WHERE w.user_id = $1
	GROUP BY w.id, w.balance;
	`
	args := []any{userId}

	var dashboard model.UserDashboardInformation
	if err := ur.db.QueryRow(ctx, sql, args...).Scan(&dashboard.Balance, &dashboard.Income, &dashboard.Expense); err != nil {
		return model.UserDashboardInformation{}, err
	}
	return dashboard, nil
}
