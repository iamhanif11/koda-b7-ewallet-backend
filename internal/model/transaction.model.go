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

type TransferDetail struct {
	ID            int       `db:"id"`
	TransactionID int       `db:"transaction_id"`
	ReceiverID    int       `db:"receiver_id"`
	Notes         *string   `db:"notes"`
	CreatedAt     time.Time `db:"created_at"`
}

type Wallet struct {
	ID        int        `db:"id" json:"id"`
	UserID    int        `db:"user_id" json:"user_id"`
	Balance   int        `db:"balance" json:"balance"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt *time.Time `db:"update_at" json:"updated_at"`
}

type TransactionHistory struct {
	Id              int
	Type            string
	Amount          int
	Status          string
	CreatedAt       time.Time
	Fullname        *string
	Picture         *string
	Phone           *string
	PaymentMethodId *int
}

// type TransactionDetail struct {
// 	ID              int     `json:"id"`
// 	UserID          int     `json:"user_id"`
// 	Amount          int     `json:"amount"`
// 	TransactionType string  `json:"transaction_type"`
// 	Status          string  `json:"status"`
// 	ReceiverID      int     `json:"receiver_id,omitempty"`
// 	ReceiverName    string  `json:"receiver_name,omitempty"`
// 	ReceiverPhone   string  `json:"receiver_phone,omitempty"`
// 	Notes           *string `json:"notes,omitempty"`
// 	CreatedAt       string  `json:"created_at"`
// }
