package dto

import "time"

type NewUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	Id        int        `json:"id"`
	Email     string     `json:"email"`
	Password  string     `json:"password"`
	Fullname  string     `json:"fullname"`
	Pin       string     `json:"pin"`
	Picture   string     `json:"picture"`
	Phone     string     `json:"phone"`
	CreatedAt *time.Time `json:"created_at"`
	UpdateAt  *time.Time `json:"updated_at"`
}
