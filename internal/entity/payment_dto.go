package entity

// CreatePaymentRequest represents the request to create a payment
type CreatePaymentRequest struct {
	BookingID uint64  `json:"bookingId" binding:"required"`
	Amount    float64 `json:"amount" binding:"required,min=0"`
	Method    string  `json:"method" binding:"required,oneof=paypal"`
}

// PayPalWebhookRequest represents PayPal webhook payload
type PayPalWebhookRequest struct {
	ID           string                 `json:"id"`
	EventType    string                 `json:"event_type"`
	CreateTime   string                 `json:"create_time"`
	ResourceType string                 `json:"resource_type"`
	Resource     map[string]interface{} `json:"resource"`
}

// PaymentResponse represents payment response
type PaymentResponse struct {
	ID        uint64   `json:"id"`
	BookingID uint64   `json:"bookingId"`
	PayerID   uint64   `json:"payerId"`
	Amount    *float64 `json:"amount,omitempty"`
	Method    *string  `json:"method,omitempty"`
	Status    string   `json:"status"`
	CreatedAt string   `json:"createdAt"`
	UpdatedAt string   `json:"updatedAt"`
	Booking   *Booking `json:"booking,omitempty"`
	Payer     *User    `json:"payer,omitempty"`
	PayPalURL *string  `json:"paypalUrl,omitempty"` // For redirect to PayPal
}

// PaymentListRequest represents filters for listing payments
type PaymentListRequest struct {
	PayerID *uint64 `json:"payerId" form:"payerId"`
	Status  *string `json:"status" form:"status" binding:"omitempty,oneof=pending success failed"`
	Page    int     `json:"page" form:"page"`
	Limit   int     `json:"limit" form:"limit"`
}

// PaymentListResponse holds paginated list of payments
type PaymentListResponse struct {
	Payments   []*PaymentResponse `json:"payments"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	Limit      int                `json:"limit"`
	TotalPages int                `json:"totalPages"`
}

// PayPalOrderRequest represents PayPal order creation request
type PayPalOrderRequest struct {
	Intent             string               `json:"intent"`
	PurchaseUnits      []PayPalPurchaseUnit `json:"purchase_units"`
	ApplicationContext PayPalAppContext     `json:"application_context"`
}

type PayPalPurchaseUnit struct {
	ReferenceID string       `json:"reference_id"`
	Amount      PayPalAmount `json:"amount"`
	Description string       `json:"description"`
}

type PayPalAmount struct {
	CurrencyCode string `json:"currency_code"`
	Value        string `json:"value"`
}

type PayPalAppContext struct {
	ReturnURL string `json:"return_url"`
	CancelURL string `json:"cancel_url"`
}

// PayPalOrderResponse represents PayPal order creation response
type PayPalOrderResponse struct {
	ID     string       `json:"id"`
	Status string       `json:"status"`
	Links  []PayPalLink `json:"links"`
}

type PayPalLink struct {
	Href   string `json:"href"`
	Rel    string `json:"rel"`
	Method string `json:"method"`
}
