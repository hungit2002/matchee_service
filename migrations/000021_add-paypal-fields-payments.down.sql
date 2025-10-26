-- Remove PayPal fields from payments table
ALTER TABLE payments
DROP COLUMN paypal_order_id,
DROP COLUMN paypal_transaction_id;

-- Revert method enum to original
ALTER TABLE payments
MODIFY COLUMN method ENUM('momo','zalo','cash','bank_transfer') DEFAULT NULL;
