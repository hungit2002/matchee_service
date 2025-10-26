package usecase

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"matchee/services/internal/entity"
	"matchee/services/internal/repository"
)

type PaymentUsecase interface {
	CreatePayment(userID uint64, req *entity.CreatePaymentRequest) (*entity.PaymentResponse, error)
	GetPayments(userID uint64, req *entity.PaymentListRequest) (*entity.PaymentListResponse, error)
	GetPaymentByID(paymentID uint64) (*entity.PaymentResponse, error)
	HandlePayPalWebhook(req *entity.PayPalWebhookRequest) error
}

type paymentUsecase struct {
	paymentRepo    repository.PaymentRepository
	bookingRepo    repository.BookingRepository
	userRepo       repository.UserRepository
	paypalClientID string
	paypalSecret   string
	paypalBaseURL  string
}

func NewPaymentUsecase(paymentRepo repository.PaymentRepository, bookingRepo repository.BookingRepository, userRepo repository.UserRepository, paypalClientID, paypalSecret, paypalBaseURL string) PaymentUsecase {
	return &paymentUsecase{
		paymentRepo:    paymentRepo,
		bookingRepo:    bookingRepo,
		userRepo:       userRepo,
		paypalClientID: paypalClientID,
		paypalSecret:   paypalSecret,
		paypalBaseURL:  paypalBaseURL,
	}
}

func (u *paymentUsecase) CreatePayment(userID uint64, req *entity.CreatePaymentRequest) (*entity.PaymentResponse, error) {
	// Verify booking exists
	_, err := u.bookingRepo.GetBookingByID(req.BookingID)
	if err != nil {
		return nil, errors.New("booking not found")
	}

	// Check if payment already exists for this booking
	existingPayment, _ := u.paymentRepo.GetPaymentByBooking(req.BookingID)
	if existingPayment != nil {
		return nil, errors.New("payment already exists for this booking")
	}

	// Create payment record
	payment := &entity.Payment{
		BookingID: req.BookingID,
		PayerID:   userID,
		Amount:    &req.Amount,
		Method:    &req.Method,
		Status:    "pending",
	}

	if err := u.paymentRepo.CreatePayment(payment); err != nil {
		return nil, err
	}

	// Create PayPal order
	paypalOrder, err := u.createPayPalOrder(payment)
	if err != nil {
		return nil, fmt.Errorf("failed to create PayPal order: %v", err)
	}

	// Get the created payment with relationships
	createdPayment, err := u.paymentRepo.GetPaymentByID(payment.ID)
	if err != nil {
		return nil, err
	}

	response := u.toPaymentResponse(createdPayment)

	// Add PayPal approval URL
	for _, link := range paypalOrder.Links {
		if link.Rel == "approve" {
			response.PayPalURL = &link.Href
			break
		}
	}

	return response, nil
}

func (u *paymentUsecase) GetPayments(userID uint64, req *entity.PaymentListRequest) (*entity.PaymentListResponse, error) {
	page := req.Page
	if page <= 0 {
		page = 1
	}
	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}

	list, total, err := u.paymentRepo.GetPaymentsByPayer(userID, req.Status, page, limit)
	if err != nil {
		return nil, err
	}

	res := make([]*entity.PaymentResponse, len(list))
	for i, p := range list {
		res[i] = u.toPaymentResponse(p)
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))
	return &entity.PaymentListResponse{Payments: res, Total: total, Page: page, Limit: limit, TotalPages: totalPages}, nil
}

func (u *paymentUsecase) GetPaymentByID(paymentID uint64) (*entity.PaymentResponse, error) {
	p, err := u.paymentRepo.GetPaymentByID(paymentID)
	if err != nil {
		return nil, errors.New("payment not found")
	}
	return u.toPaymentResponse(p), nil
}

func (u *paymentUsecase) HandlePayPalWebhook(req *entity.PayPalWebhookRequest) error {
	// Handle different PayPal webhook events
	switch req.EventType {
	case "PAYMENT.CAPTURE.COMPLETED":
		return u.handlePaymentCompleted(req)
	case "PAYMENT.CAPTURE.DENIED":
		return u.handlePaymentFailed(req)
	default:
		// Log unhandled event types
		return nil
	}
}

func (u *paymentUsecase) createPayPalOrder(payment *entity.Payment) (*entity.PayPalOrderResponse, error) {
	// Get PayPal access token
	accessToken, err := u.getPayPalAccessToken()
	if err != nil {
		return nil, err
	}

	// Create order request
	orderReq := &entity.PayPalOrderRequest{
		Intent: "CAPTURE",
		PurchaseUnits: []entity.PayPalPurchaseUnit{
			{
				ReferenceID: fmt.Sprintf("booking_%d", payment.BookingID),
				Amount: entity.PayPalAmount{
					CurrencyCode: "USD",
					Value:        fmt.Sprintf("%.2f", *payment.Amount),
				},
				Description: fmt.Sprintf("Payment for booking #%d", payment.BookingID),
			},
		},
		ApplicationContext: entity.PayPalAppContext{
			ReturnURL: "https://your-app.com/payment/success",
			CancelURL: "https://your-app.com/payment/cancel",
		},
	}

	// Make request to PayPal
	jsonData, _ := json.Marshal(orderReq)
	req, _ := http.NewRequest("POST", u.paypalBaseURL+"/v2/checkout/orders", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		return nil, errors.New("failed to create PayPal order")
	}

	var orderResp entity.PayPalOrderResponse
	if err := json.NewDecoder(resp.Body).Decode(&orderResp); err != nil {
		return nil, err
	}
	return &orderResp, nil
}

func (u *paymentUsecase) getPayPalAccessToken() (string, error) {
	req, _ := http.NewRequest("POST", u.paypalBaseURL+"/v1/oauth2/token", nil)
	req.SetBasicAuth(u.paypalClientID, u.paypalSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Body = io.NopCloser(bytes.NewBufferString("grant_type=client_credentials"))

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", err
	}
	return tokenResp.AccessToken, nil
}

func (u *paymentUsecase) handlePaymentCompleted(req *entity.PayPalWebhookRequest) error {
	// Extract payment ID from webhook data
	// This is a simplified implementation - you'd need to parse the actual PayPal webhook structure
	// and match it to your payment records

	// For now, we'll update the payment status to success
	// In a real implementation, you'd need to:
	// 1. Extract the PayPal transaction ID
	// 2. Find the corresponding payment record
	// 3. Update the status

	return nil
}

func (u *paymentUsecase) handlePaymentFailed(req *entity.PayPalWebhookRequest) error {
	// Similar to handlePaymentCompleted but set status to failed
	return nil
}

func (u *paymentUsecase) toPaymentResponse(p *entity.Payment) *entity.PaymentResponse {
	return &entity.PaymentResponse{
		ID:        p.ID,
		BookingID: p.BookingID,
		PayerID:   p.PayerID,
		Amount:    p.Amount,
		Method:    p.Method,
		Status:    p.Status,
		CreatedAt: p.CreatedAt.Format(time.RFC3339),
		UpdatedAt: p.UpdatedAt.Format(time.RFC3339),
		Booking:   p.Booking,
		Payer:     p.Payer,
	}
}
