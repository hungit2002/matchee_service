-- Insert default roles
INSERT INTO roles (name, created_at, updated_at) VALUES 
('player', NOW(), NOW()),
('admin', NOW(), NOW()),
('venue_owner', NOW(), NOW())
ON DUPLICATE KEY UPDATE updated_at = NOW();
