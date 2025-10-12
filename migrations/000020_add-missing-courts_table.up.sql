-- Add missing columns to courts table
ALTER TABLE courts 
ADD COLUMN description TEXT NULL,
ADD COLUMN court_type ENUM('indoor','outdoor') DEFAULT 'outdoor' NOT NULL,
ADD COLUMN surface ENUM('hard','clay','grass','synthetic') DEFAULT 'hard' NOT NULL,
ADD COLUMN is_active BOOLEAN DEFAULT TRUE NOT NULL;

-- Remove unused column
ALTER TABLE courts 
DROP COLUMN court_number;
