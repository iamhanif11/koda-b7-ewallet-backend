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
