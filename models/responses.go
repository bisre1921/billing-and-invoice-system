package models

type LoginCredentials struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

type GenericResponse struct {
	Message string      `json:"message"`
	ID      interface{} `json:"id,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
