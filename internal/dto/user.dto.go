package dto

type UserProfileRes struct {
	Fullname *string `json:"fullname"`
	Email    string  `json:"email"`
	Picture  *string `json:"picture"`
}

type UserCheckPinRes struct {
	IsValid bool `json:"isvalid"`
}
type UserCheckPinReq struct {
	Pin string `json:"pin" binding:"required, len=6"`
}

type UserUpdateProfileReq struct {
	Fullname *string `json:"fullname"`
	Phone    *string `json:"phone"`
	Picture  *string `json:"picture"`
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
	Pin string `json:"pin" binding:"required,len=6"`
}
