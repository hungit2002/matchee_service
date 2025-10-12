CREATE TABLE match_groups (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    match_post_id BIGINT NULL,
    venue_id BIGINT NULL,
    court_id BIGINT NULL,
    scheduled_time TIMESTAMP NULL,
    total_price DECIMAL(10,2) NULL,
    status ENUM('pending','confirmed','completed','cancelled') DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    FOREIGN KEY (match_post_id) REFERENCES match_posts(id) ON DELETE SET NULL,
    FOREIGN KEY (venue_id) REFERENCES venues(id) ON DELETE SET NULL,
    FOREIGN KEY (court_id) REFERENCES courts(id) ON DELETE SET NULL
);
