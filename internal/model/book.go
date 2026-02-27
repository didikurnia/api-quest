package model

// Book represents a book entity stored in memory.
type Book struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   int    `json:"year"`
}

// CreateBookRequest is the payload for creating a new book.
type CreateBookRequest struct {
	Title  string `json:"title" binding:"required"`
	Author string `json:"author" binding:"required"`
	Year   int    `json:"year" binding:"required"`
}

// UpdateBookRequest is the payload for updating an existing book.
type UpdateBookRequest struct {
	Title  string `json:"title" binding:"required"`
	Author string `json:"author" binding:"required"`
	Year   int    `json:"year" binding:"required"`
}

// AuthRequest is the payload for obtaining a JWT token.
type AuthRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// ErrorResponse is a standard error shape.
type ErrorResponse struct {
	Error string `json:"error"`
}

// TokenResponse is the response for a successful auth request.
type TokenResponse struct {
	Token string `json:"token"`
}

// PaginatedResponse wraps a list of books with pagination metadata.
type PaginatedResponse struct {
	Data  []Book `json:"data"`
	Page  int    `json:"page"`
	Limit int    `json:"limit"`
	Total int    `json:"total"`
}
