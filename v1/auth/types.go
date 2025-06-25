package auth

type RegisterRequest struct {
	Email     string
	Password  string
	FirstName string
	LastName  string
}

type LoginRequest struct {
	Email    string
	Password string
}

type AuthResponse struct {
	Token string `json:"token"`
}
