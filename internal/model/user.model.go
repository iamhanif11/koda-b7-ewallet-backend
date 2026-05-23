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

type UserDashboardInformation struct {
	Balance int `db:"balance"`
	Income  int `db:"income"`
	Expense int `db:"expense"`
}

type UserTransactionReport struct {
	Date    time.Time `json:"date"`
	Income  int       `json:"income"`
	Expense int       `json:"expense"`
}
