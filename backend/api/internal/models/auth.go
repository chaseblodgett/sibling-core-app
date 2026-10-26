package models

// SignUpRequest is the JSON body for POST /auth/signup.
type SignUpRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

// LoginRequest is the JSON body for POST /auth/login.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// GoogleSignInRequest is the JSON body for POST /auth/google. IDToken is
// the credential returned by Google Identity Services on the frontend.
type GoogleSignInRequest struct {
	IDToken string `json:"id_token"`
}

// SignUpInput carries the fields needed to register a new user.
type SignUpInput struct {
	Email    string
	Password string
	Name     string
}

// LoginInput carries credentials for authenticating an existing user.
type LoginInput struct {
	Email    string
	Password string
}

// AuthResponse is returned after a successful sign up or login.
type AuthResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}
