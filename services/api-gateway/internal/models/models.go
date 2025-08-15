package models

type LoginResponse struct {
    Token string `json:"token"`
    User  User   `json:"user"`
}

type User struct {
    ID    string `json:"id"`
    Email string `json:"email"`
    Role  string `json:"role"`
}