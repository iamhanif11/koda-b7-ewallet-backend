package model

import "time"

type User struct {
	Id        int        `db:"id"`
	Email     string     `db:"email"`
	Password  string     `db:"password"`
	Fullname  *string    `db:"fullname"`
	Pin       *string    `db:"pin"`
	Picture   *string    `db:"picture"`
	Phone     *string    `db:"phone"`
	CreatedAt *time.Time `db:"created_at"`
	UpdateAt  *time.Time `db:"updated_at"`
}

// type Wallet struct {
// 	Id        int        `db:"id"`
// 	User_id   int        `db:"user_id"`
// 	Balance   int        `db:"balance"`
// 	CreatedAt *time.Time `db:"created_at"`
// 	UpdateAt  *time.Time `db:"updated_at"`
// }

// type Payment_method struct {
// 	Id   int    `db:"id"`
// 	Name string `db:"name"`
// }

// type Transaction struct {
// 	Id               int        `db:"id"`
// 	User_id          int        `db:"user_id"`
// 	Amount           int        `db:"balance"`
// 	Transaction_type string     `db:"transaction_type"`
// 	Status           string     `db:"status"`
// 	CreatedAt        *time.Time `db:"created_at"`
// 	UpdatedAt        *time.Time `db:"updated_at"`
// }

// type Transfer_detail struct {
// 	Id             int        `db:"id"`
// 	Transaction_id int        `db:"transaction_id"`
// 	Receiver_id    int        `db:"receiver_id"`
// 	Notes          string     `db:"notes"`
// 	CreatedAt      *time.Time `db:"created_at"`
// }

// type Topup_detail struct {
// 	Id                int        `db:"id"`
// 	Transaction_id    int        `db:"transaction_id"`
// 	Payment_method_id int        `db:"payment_method_id"`
// 	ServiceFee        int        `db:"service_fee"`
// 	TaxAmount         int        `db:"tax_amount"`
// 	SubTotal          int        `db:"sub_total"`
// 	CreatedAt         *time.Time `db:"created_at"`
// }
