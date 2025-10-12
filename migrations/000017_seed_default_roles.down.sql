-- Remove default roles
DELETE FROM roles WHERE name IN ('player', 'admin', 'venue_owner');
