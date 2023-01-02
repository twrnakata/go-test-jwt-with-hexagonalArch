package controllers

type UserController interface {
	Signup(SignupRequest) (*SignupResponse, error)
	Login(LoginRequest) (*LoginResponse, error)
	View(int) (*LoginResponse, error)
}

type SignupRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type SignupResponse struct {
	Username string `json:"username"`
	Role     string `json:"role"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}
