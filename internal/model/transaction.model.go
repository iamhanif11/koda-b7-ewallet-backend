package model

import "time"

type Transaction struct {
	Id              int        `db:"id"`
	UserId          int        `db:"user_id"`
	Amount          int        `db:"amount"`
	TransactionType string     `db:"transaction_type"`
	Status          string     `db:"status"`
	CreatedAt       time.Time  `db:"created_at"`
	UpadatedAt      *time.Time `db:"updated_at"`
}

type Receiver struct {
	Id       int     `db:"id"`
	Picture  *string `db:"picture"`
	Receiver string  `db:"fullname"`
	Phone    string  `db:"phone"`
}
