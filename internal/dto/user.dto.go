package dto

import "mime/multipart"

type UserProfileRes struct {
	Fullname *string `json:"fullname"`
	Email    string  `json:"email"`
	Picture  *string `json:"picture"`
}

type UserCheckPinRes struct {
	IsValid bool `json:"isvalid"`
}
type UserCheckPinReq struct {
	Pin string `json:"pin" binding:"required,len=6"`
}

type UserUpdateProfileReq struct {
	Fullname *string               `form:"fullname"`
	Phone    *string               `form:"phone"`
	Picture  *multipart.FileHeader `form:"picture" binding:"omitempty"`
}

type UserUpdateProfilRes struct {
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Picture  string `json:"picture"`
}

type UserUpdatePasswordReq struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
}

type UserUpdatePasswordRes struct {
	Fullname *string `json:"fullname"`
	Email    string  `json:"email"`
	Phone    *string `json:"phone"`
	Picture  *string `json:"picture"`
}

type UserUpdatePinReq struct {
	Pin string `json:"pin" binding:"required,numeric,len=6"`
}

type UserDashboardInformationRes struct {
	Balance int `json:"balance"`
	Income  int `json:"income"`
	Expense int `json:"expense"`
}

type UserTransactionReportRes struct {
	Date    string `json:"date"`
	Day     string `json:"day"`
	Income  int    `json:"income"`
	Expense int    `json:"expense"`
}
