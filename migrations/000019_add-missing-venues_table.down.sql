-- Rollback changes to venues table
-- Rename columns back to original names
ALTER TABLE venues 
CHANGE COLUMN phone contact_phone VARCHAR(20) NULL,
CHANGE COLUMN price_per_hour base_price DECIMAL(10,2) NULL;

-- Remove added columns
ALTER TABLE venues 
DROP COLUMN description,
DROP COLUMN email,
DROP COLUMN website,
DROP COLUMN is_active;
