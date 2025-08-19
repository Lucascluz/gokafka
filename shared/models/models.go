package models

type UserData struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string   `json:"token"`
	Data  UserData `json:"data"`
}

type RegisterRequest struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type GetProfileRequest struct {
	ID string `json:"id"`
}

type GetProfileResponse struct {
	Status string   `json:"status"`
	Data   UserData `json:"data"`
}

type ListProfileResponse struct {
	Status string     `json:"status"`
	Data   []UserData `json:"data"`
}
