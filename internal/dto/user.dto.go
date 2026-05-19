package dto

type UserProfileRes struct {
	Fullname *string `json:"fullname"`
	Email    string  `json:"email"`
	Picture  *string `json:"picture"`
}
