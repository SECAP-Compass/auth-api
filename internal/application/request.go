package application

type UserRegisterRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	Authority string `json:"authority"`
}

type UserLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
