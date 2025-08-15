package models

type Request struct {
	Type          string `json:"type"` // e.g., "register", "login", etc.
	CorrelationID string `json:"correlation_id"`
	ReplyTo       string `json:"reply_to"`
	Payload       string `json:"payload"` // You can use a string or json.RawMessage for more complex payloads
}

type Response struct {
    CorrelationID string `json:"correlation_id"`
    Success       bool   `json:"success"`
    Data          string `json:"data"`
    Error         string `json:"error,omitempty"`
}