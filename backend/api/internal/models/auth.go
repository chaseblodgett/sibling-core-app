package models

// SignUpRequest is the JSON body for POST /auth/signup.
type SignUpRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
}

// LoginRequest is the JSON body for POST /auth/login.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// SignUpInput carries the fields needed to register a new user.
type SignUpInput struct {
	Email    string
	Password string
	FullName string
}

// LoginInput carries credentials for authenticating an existing user.
type LoginInput struct {
	Email    string
	Password string
}

// AuthResult is returned after a successful sign up or login.
type AuthResult struct {
	UserID       string
	AccessToken  string
	RefreshToken string
}
