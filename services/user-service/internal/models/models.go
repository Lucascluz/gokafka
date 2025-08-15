package models

type User struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Role      string `json:"role"`
}

type Session struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	ExpiresAt string `json:"expires_at"`
	Role	  string `json:"role"`
}

type Request struct {
	Type          string `json:"type"` // e.g., "register", "login", etc.
	CorrelationID string `json:"correlation_id"`
	ReplyTo       string `json:"reply_to"`
	Payload       string `json:"payload"` // You can use a string or json.RawMessage for more complex payloads
}

type Response struct {
	CorrelationID string `json:"correlation_id"`
	Data          string `json:"data"`
}
