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

type TransferRequest struct {
	ReceiverId int    `json:"receiver_id" validate:"required,gt=0"`
	Amount     int    `json:"amount" validate:"required,gt=0"`
	Notes      string `json:"notes" validate:"max=255"`
}

type TransferResponse struct {
	ID              int    `json:"id"`
	Amount          int    `json:"amount"`
	ReceiverID      int    `json:"receiver_id"`
	ReceiverName    string `json:"receiver_name"`
	ReceiverPhone   string `json:"receiver_phone"`
	Notes           string `json:"notes"`
	Status          string `json:"status"`
	TransactionType string `json:"transaction_type"`
	CreatedAt       string `json:"created_at"`
}

type TopUpRequest struct {
	Amount          int `json:"amount"`
	PaymentMethodId int `json:"payment_method_id"`
}
