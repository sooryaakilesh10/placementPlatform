package presenter

type UserRequest struct {
	UserName  string `json:"user_name"`
	Email     string `json:"email"`
	Pass      string `json:"pass"`
	Role      string `json:"role"`
	CreatedBy uint64 `json:"created_by"`
}

type UserResponse struct {
	UserID   string `json:"user_id"`
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type LoginRequest struct {
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	Pass     string `json:"pass"`
}

type LoginResponse struct {
	JWTToken string `json:"jwt_token"`
}
