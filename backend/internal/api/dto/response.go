package dto

// ErrorResponse representa una respuesta de error
type ErrorResponse struct {
	Error    string `json:"error"`
	Message  string `json:"message"`
	Code     string `json:"code,omitempty"`
}

// SuccessResponse representa una respuesta exitosa generica
type SuccessResponse struct {
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
}

// PaginatedResponse representa una respuesta paginada
type PaginatedResponse struct {
	Data          interface{} `json:"data"`
	Total         int         `json:"total"`
	Limit         int         `json:"limit"`
	Offset        int         `json:"offset"`
	TotalPages    int         `json:"total_pages"`
}