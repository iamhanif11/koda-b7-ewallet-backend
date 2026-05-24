package repository

import (
	"context"

	"github.com/iamhanif11/ewallet-backend/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type DBTX interface {
	Query(context.Context, string, ...any) (pgx.Rows, error)
	QueryRow(context.Context, string, ...any) pgx.Row
	Exec(context.Context, string, ...any) (pgconn.CommandTag, error)
}

type TransactionRepository struct{}

func NewTransactionRepository() *TransactionRepository {
	return &TransactionRepository{}
}

func (tr *TransactionRepository) FindReceivers(ctx context.Context, dbtx DBTX, userId int, search string, limit, offset int) ([]model.Receiver, error) {
	sql := `
		SELECT id, picture, fullname AS receiver, phone
		FROM users
		WHERE id = $1 AND
			(
				fullname ILIKE $2 || '%'
				OR phone ILIKE $2 || '%'
			)
		ORDER BY fullname ASC
		LIMIT $3
		OFFSET $4
	`
	args := []any{userId, search, limit, offset}

	rows, err := dbtx.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	receivers := []model.Receiver{}
	for rows.Next() {
		var receiver model.Receiver
		if err := rows.Scan(&receiver.Id, &receiver.Picture, &receiver.Receiver, &receiver.Phone); err != nil {
			return nil, err
		}
		receivers = append(receivers, receiver)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return receivers, nil
}
