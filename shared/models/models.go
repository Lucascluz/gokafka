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

// Product-related models
type ProductData struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

type CreateProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type UpdateProductRequest struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type GetProductRequest struct {
	ID int `json:"id"`
}

type DeleteProductRequest struct {
	ID int `json:"id"`
}

type GetProductResponse struct {
	Status string      `json:"status"`
	Data   ProductData `json:"data"`
}

type ListProductResponse struct {
	Status string        `json:"status"`
	Data   []ProductData `json:"data"`
}

type ProductResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    ProductData `json:"data"`
}
