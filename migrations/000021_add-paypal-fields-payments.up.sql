-- Add PayPal fields to payments table
ALTER TABLE payments
ADD COLUMN paypal_order_id VARCHAR(255) NULL,
ADD COLUMN paypal_transaction_id VARCHAR(255) NULL;

-- Update method enum to include paypal
ALTER TABLE payments
MODIFY COLUMN method ENUM('paypal','momo','zalo','cash','bank_transfer') DEFAULT 'paypal';
