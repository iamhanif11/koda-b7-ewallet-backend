package dto

import "time"

type NewUser struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"min=6"`
}

type User struct {
	Id        int        `json:"id"`
	Email     string     `json:"email"`
	Password  string     `json:"password,omitempty"`
	Fullname  string     `json:"fullname,omitempty"`
	Pin       string     `json:"pin,omitempty"`
	Picture   string     `json:"picture,omitempty"`
	Phone     string     `json:"phone,omitempty"`
	CreatedAt *time.Time `json:"created_at"`
	UpdateAt  *time.Time `json:"updated_at"`
}
