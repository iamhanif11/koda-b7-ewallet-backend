package dto

type ReceiverResponse struct {
	Id       int     `json:"id"`
	Picture  *string `json:"picture"`
	Receiver string  `json:"reciver"`
	Phone    string  `json:"phone"`
}

type PaginationResponse struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

type ReceiverListResponse struct {
	Items []ReceiverResponse `json:"items"`
	Pages PaginationResponse `json:"pages"`
}
