package repository

import (
	"context"
	"errors"

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
		WHERE id != $1 AND
			(
				fullname ILIKE	'%' || $2 || '%'
				OR phone ILIKE	'%' || $2 || '%'
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

func (tr *TransactionRepository) Transfer(ctx context.Context, dbtx DBTX, senderId, receiverId int, amount int, notes string) error {
	var balance int
	checkBalanceSql := `
		SELECT balance
		FROM wallet
		WHERE user_id = $1
		FOR UPDATE
	`

	err := dbtx.QueryRow(ctx, checkBalanceSql, senderId).Scan(&balance)

	if err != nil {
		return err
	}

	if balance < amount {
		return errors.New("Insufficient balance")
	}

	//tx
	var transactionId int

	insertTxSql := `
		INSERT INTO transactions (
			user_id, amount, transaction_type, status
		) VALUES ($1, $2, 'transfer out', 'completed')
		RETURNING id
	`
	err = dbtx.QueryRow(ctx, insertTxSql, senderId, amount).Scan(&transactionId)

	if err != nil {
		return err
	}

	//tx detail
	insertTxDetailSql := `
		INSERT INTO transfer_detail (
			transaction_id,
			receiver_id,
			notes
		)
		VALUES($1, $2, $3)
	`

	_, err = dbtx.Exec(ctx, insertTxDetailSql, transactionId, receiverId, notes)

	if err != nil {
		return err
	}

	//reduce sender balance
	reduceBalanceSql := `
		UPDATE wallet
		SET
			balance = balance - $1,
			update_at = CURRENT_TIMESTAMP
		WHERE user_id = $2
	`

	_, err = dbtx.Exec(ctx, reduceBalanceSql, amount, senderId)

	if err != nil {
		return err
	}

	//add balance receive
	addBalanceSql := `
		UPDATE wallet
		SET
			balance = balance + $1,
			update_at = CURRENT_TIMESTAMP
		WHERE user_id = $2
	`

	_, err = dbtx.Exec(ctx, addBalanceSql, amount, receiverId)

	if err != nil {
		return err
	}

	return nil
}

func (tr *TransactionRepository) TopUp(ctx context.Context, dbtx DBTX, userId, amount, paymentMethodId int) error {
	var paymentMethodExist bool

	checkPaymentMethodSql := `
		SELECT EXISTS (
			SELECT 1
			FROM payment_method
			WHERE id = $1
		)
	`

	err := dbtx.QueryRow(ctx, checkPaymentMethodSql, paymentMethodId).Scan(&paymentMethodExist)

	if err != nil {
		return err
	}

	if !paymentMethodExist {
		return errors.New("Payment method not found")
	}

	//calc
	serviceFee := 0
	taxAmount := amount * 10 / 100
	subTotal := amount + serviceFee + taxAmount

	//insert topup
	var transactionId int

	insertTransactionSql := `
		INSERT INTO transactions (
			user_id, amount, transaction_type, status
		)
		VALUES ($1, $2, 'top up', 'completed')
		RETURNING id
	`

	err = dbtx.QueryRow(
		ctx, insertTransactionSql, userId, amount,
	).Scan(&transactionId)

	if err != nil {
		return err
	}

	//insert top up
	topUpDetailSql := `
		INSERT INTO topup_detail(
			transaction_id, payment_method_id, service_fee, tax_amount,sub_total
		)
		VALUES($1, $2, $3, $4, $5)
	`

	_, err = dbtx.Exec(
		ctx, topUpDetailSql, transactionId, paymentMethodId, serviceFee, taxAmount, subTotal,
	)

	if err != nil {
		return err
	}

	updateWalletSql := `
		UPDATE wallet
		SET
			balance = balance + $1,
			update_at = CURRENT_TIMESTAMP
		WHERE user_id = $2
	`

	result, err := dbtx.Exec(ctx, updateWalletSql, amount, userId)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("wallet not found")
	}

	return nil
}

func (tr *TransactionRepository) GetTransactionHistoryById(ctx context.Context, dbtx DBTX, userId int, search string, limit int, offset int) ([]model.TransactionHistory, error) {
	sql := `
		SELECT id, transaction_type, amount, status, created_at
		FROM transactions
		WHERE user_id = $1
		AND (
			transaction_type ILIKE '%' || $2 || '%'
			OR status ILIKE '%' || $2 || '%'
		)
		ORDER BY created_at DESC
		LIMIT $3
		OFFSET $4
	`

	args := []any{userId, search, limit, offset}
	rows, err := dbtx.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var histories []model.TransactionHistory

	for rows.Next() {
		var history model.TransactionHistory

		err := rows.Scan(&history.Id, &history.Type, &history.Amount, &history.Status, &history.CreatedAt)

		if err != nil {
			return nil, err
		}

		histories = append(histories, history)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return histories, nil

}
