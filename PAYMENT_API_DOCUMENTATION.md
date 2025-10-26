# Payment API Documentation

## Overview
Payment management APIs with PayPal integration for court booking payments.

## Endpoints

### 1. Create Payment Request
**POST** `/api/v1/payments`

Creates a payment request for a booking with PayPal integration.

**Headers:**
- `Authorization: Bearer {access_token}`
- `Content-Type: application/json`

**Request Body:**
```json
{
  "bookingId": 1,
  "amount": 400000,
  "method": "paypal"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Payment request created successfully",
  "data": {
    "id": 1,
    "bookingId": 1,
    "payerId": 1,
    "amount": 400000,
    "method": "paypal",
    "status": "pending",
    "createdAt": "2024-01-25T10:00:00Z",
    "updatedAt": "2024-01-25T10:00:00Z",
    "paypalUrl": "https://www.sandbox.paypal.com/checkoutnow?token=5O190127TN364715T",
    "booking": {
      "id": 1,
      "venueId": 1,
      "courtId": 1,
      "startTime": "2024-01-25T18:00:00Z",
      "endTime": "2024-01-25T20:00:00Z",
      "totalPrice": 400000,
      "status": "reserved"
    },
    "payer": {
      "id": 1,
      "email": "user@example.com",
      "fullName": "John Doe"
    }
  }
}
```

### 2. Get Payment History
**GET** `/api/v1/payments`

Retrieves payment history for the current user with optional filters.

**Headers:**
- `Authorization: Bearer {access_token}`

**Query Parameters:**
- `status` (optional): Filter by payment status (`pending`, `success`, `failed`)
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 20)

**Response:**
```json
{
  "success": true,
  "message": "Payments retrieved successfully",
  "data": {
    "payments": [
      {
        "id": 1,
        "bookingId": 1,
        "payerId": 1,
        "amount": 400000,
        "method": "paypal",
        "status": "success",
        "createdAt": "2024-01-25T10:00:00Z",
        "updatedAt": "2024-01-25T10:05:00Z",
        "booking": {
          "id": 1,
          "venueId": 1,
          "courtId": 1,
          "startTime": "2024-01-25T18:00:00Z",
          "endTime": "2024-01-25T20:00:00Z",
          "totalPrice": 400000,
          "status": "paid"
        },
        "payer": {
          "id": 1,
          "email": "user@example.com",
          "fullName": "John Doe"
        }
      }
    ],
    "total": 1,
    "page": 1,
    "limit": 20,
    "totalPages": 1
  }
}
```

### 3. Get Payment Details
**GET** `/api/v1/payments/{id}`

Retrieves detailed payment information by ID.

**Headers:**
- `Authorization: Bearer {access_token}`

**Response:**
```json
{
  "success": true,
  "message": "Payment retrieved successfully",
  "data": {
    "id": 1,
    "bookingId": 1,
    "payerId": 1,
    "amount": 400000,
    "method": "paypal",
    "status": "success",
    "paypalOrderId": "5O190127TN364715T",
    "paypalTransactionId": "2GG279541U8729312",
    "createdAt": "2024-01-25T10:00:00Z",
    "updatedAt": "2024-01-25T10:05:00Z",
    "booking": {
      "id": 1,
      "venueId": 1,
      "courtId": 1,
      "startTime": "2024-01-25T18:00:00Z",
      "endTime": "2024-01-25T20:00:00Z",
      "totalPrice": 400000,
      "status": "paid"
    },
    "payer": {
      "id": 1,
      "email": "user@example.com",
      "fullName": "John Doe"
    }
  }
}
```

### 4. PayPal Webhook
**POST** `/api/v1/payments/webhook`

Handles PayPal webhook notifications for payment status updates.

**Headers:**
- `Content-Type: application/json`

**Request Body (PayPal Webhook):**
```json
{
  "id": "WH-2W4268056S058805L-67976317FL0907054",
  "event_type": "PAYMENT.CAPTURE.COMPLETED",
  "create_time": "2018-12-10T21:20:49.000Z",
  "resource_type": "capture",
  "resource": {
    "id": "2GG279541U8729312",
    "amount": {
      "currency_code": "USD",
      "value": "4.54"
    },
    "status": "COMPLETED"
  }
}
```

**Response:**
```json
{
  "success": true,
  "message": "Webhook processed successfully",
  "data": null
}
```

## PayPal Integration

### Configuration
The following environment variables need to be set:

```bash
PAYPAL_CLIENT_ID=your_paypal_client_id
PAYPAL_SECRET=your_paypal_secret
PAYPAL_BASE_URL=https://api.sandbox.paypal.com  # For sandbox
# PAYPAL_BASE_URL=https://api.paypal.com  # For production
```

### Payment Flow
1. User creates a payment request via `/api/v1/payments`
2. System creates PayPal order and returns approval URL
3. User is redirected to PayPal for payment
4. PayPal sends webhook notification to `/api/v1/payments/webhook`
5. System updates payment status based on webhook event

### Supported PayPal Events
- `PAYMENT.CAPTURE.COMPLETED`: Payment successful
- `PAYMENT.CAPTURE.DENIED`: Payment failed

## Error Responses

### 400 Bad Request
```json
{
  "success": false,
  "message": "Invalid request data",
  "error": "validation error details"
}
```

### 401 Unauthorized
```json
{
  "success": false,
  "message": "User not authenticated",
  "error": "Invalid or missing access token"
}
```

### 404 Not Found
```json
{
  "success": false,
  "message": "Payment not found",
  "error": "Payment with ID 999 does not exist"
}
```

## Data Models

### Payment Entity
```go
type Payment struct {
    ID        uint64     `json:"id"`
    BookingID uint64     `json:"bookingId"`
    PayerID   uint64     `json:"payerId"`
    Amount    *float64   `json:"amount,omitempty"`
    Method    *string    `json:"method,omitempty"`
    Status    string     `json:"status"`
    PayPalOrderID *string `json:"paypalOrderId,omitempty"`
    PayPalTransactionID *string `json:"paypalTransactionId,omitempty"`
    CreatedAt string     `json:"createdAt"`
    UpdatedAt string     `json:"updatedAt"`
    Booking   *Booking   `json:"booking,omitempty"`
    Payer     *User      `json:"payer,omitempty"`
}
```

### Payment Status Values
- `pending`: Payment request created, awaiting PayPal processing
- `success`: Payment completed successfully
- `failed`: Payment failed or was denied

### Payment Method Values
- `paypal`: PayPal payment
- `momo`: MoMo wallet (legacy)
- `zalo`: ZaloPay (legacy)
- `cash`: Cash payment (legacy)
- `bank_transfer`: Bank transfer (legacy)
