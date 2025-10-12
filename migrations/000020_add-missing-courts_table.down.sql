-- Rollback changes to courts table
-- Add back the removed column
ALTER TABLE courts 
ADD COLUMN court_number INT NULL;

-- Remove added columns
ALTER TABLE courts 
DROP COLUMN description,
DROP COLUMN court_type,
DROP COLUMN surface,
DROP COLUMN is_active;
