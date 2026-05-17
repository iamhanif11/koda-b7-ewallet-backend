package dto

type Response struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Success bool   `json:"isSucces"`
	Error   string `json:"error,omitempty"`
}
