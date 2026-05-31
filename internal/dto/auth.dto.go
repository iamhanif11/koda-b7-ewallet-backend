package dto

import "time"

type NewUser struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"min=6"`
}

type Login struct {
	Email    string `json:"email" binding:"required" example:"kentang@gmail.com"`
	Password string `json:"password" binding:"required,min=6" example:"123456"`
}

type User struct {
	Id        int        `json:"id"`
	Email     string     `json:"email"`
	Password  string     `json:"password"`
	Fullname  string     `json:"fullname,omitempty"`
	Pin       string     `json:"pin,omitempty"`
	Picture   string     `json:"picture,omitempty"`
	Phone     string     `json:"phone,omitempty"`
	CreatedAt *time.Time `json:"created_at"`
	UpdateAt  *time.Time `json:"updated_at"`
}

type LoginResponse struct {
	Token  string `json:"token" example:"token..."`
	HasPin bool   `json:"has_pin"`
	User   User   `json:"user"`
}

type VerifyEmailReq struct {
	Email string `json:"email"`
}

type ResetPasswordReq struct {
	Email           string `json:"email"`
	NewPassword     string `json:"new_password"`
	ConfirmPassword string `json:"confirm_password"`
}
