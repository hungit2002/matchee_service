-- Add missing columns to venues table
ALTER TABLE venues 
ADD COLUMN email VARCHAR(100) NULL,
ADD COLUMN website VARCHAR(255) NULL,
ADD COLUMN is_active BOOLEAN DEFAULT TRUE NOT NULL;

-- Rename existing columns to match entity
ALTER TABLE venues 
CHANGE COLUMN contact_phone phone VARCHAR(20) NOT NULL,
CHANGE COLUMN base_price price_per_hour DECIMAL(10,2) NOT NULL;
