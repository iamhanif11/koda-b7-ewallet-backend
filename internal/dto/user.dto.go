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
	Pin string `json:"pin"`
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
