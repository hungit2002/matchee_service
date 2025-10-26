package repository

import (
	"fmt"
	"time"

	"matchee/services/internal/entity"

	"gorm.io/gorm"
)

type AdminRepository interface {
	GetRevenueStatistics(fromDate, toDate string) ([]*entity.RevenueStatistics, float64, error)
	GetUsersWithStats(req *entity.AdminUserListRequest) ([]*entity.AdminUserResponse, int64, error)
	LockUser(userID uint64, locked bool, reason *string) error
	ApproveVenue(venueID uint64, approved bool, reason *string) error
	GetFeedbacksWithDetails(req *entity.AdminFeedbackListRequest) ([]*entity.AdminFeedbackResponse, int64, error)
}

type adminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) AdminRepository {
	return &adminRepository{db: db}
}

func (r *adminRepository) GetRevenueStatistics(fromDate, toDate string) ([]*entity.RevenueStatistics, float64, error) {
	var stats []*entity.RevenueStatistics
	var totalRevenue float64

	// Build date range query
	dateQuery := ""
	if fromDate != "" && toDate != "" {
		dateQuery = fmt.Sprintf("AND DATE(p.created_at) BETWEEN '%s' AND '%s'", fromDate, toDate)
	} else if fromDate != "" {
		dateQuery = fmt.Sprintf("AND DATE(p.created_at) >= '%s'", fromDate)
	} else if toDate != "" {
		dateQuery = fmt.Sprintf("AND DATE(p.created_at) <= '%s'", toDate)
	}

	// Get daily revenue statistics
	query := fmt.Sprintf(`
		SELECT 
			DATE(p.created_at) as date,
			COALESCE(SUM(p.amount), 0) as total_revenue,
			COUNT(DISTINCT p.booking_id) as booking_count,
			COALESCE(AVG(p.amount), 0) as average_order_value
		FROM payments p
		WHERE p.status = 'success' %s
		GROUP BY DATE(p.created_at)
		ORDER BY DATE(p.created_at) DESC
	`, dateQuery)

	err := r.db.Raw(query).Scan(&stats).Error
	if err != nil {
		return nil, 0, err
	}

	// Get total revenue
	totalQuery := fmt.Sprintf(`
		SELECT COALESCE(SUM(amount), 0) as total
		FROM payments 
		WHERE status = 'success' %s
	`, dateQuery)

	err = r.db.Raw(totalQuery).Scan(&totalRevenue).Error
	if err != nil {
		return nil, 0, err
	}

	return stats, totalRevenue, nil
}

func (r *adminRepository) GetUsersWithStats(req *entity.AdminUserListRequest) ([]*entity.AdminUserResponse, int64, error) {
	var users []*entity.AdminUserResponse
	var total int64

	// Build base query
	query := r.db.Model(&entity.User{})

	// Apply filters
	if req.Status != nil {
		if *req.Status == "locked" {
			query = query.Where("deleted_at IS NOT NULL")
		} else {
			query = query.Where("deleted_at IS NULL")
		}
	}

	if req.Role != nil {
		query = query.Where("role = ?", *req.Role)
	}

	if req.Search != nil && *req.Search != "" {
		searchTerm := "%" + *req.Search + "%"
		query = query.Where("email LIKE ? OR full_name LIKE ?", searchTerm, searchTerm)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	page := req.Page
	if page <= 0 {
		page = 1
	}
	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}
	offset := (page - 1) * limit

	// Get users with additional stats
	userQuery := `
		SELECT 
			u.id,
			u.email,
			u.full_name,
			u.phone,
			u.role,
			CASE WHEN u.deleted_at IS NULL THEN 'active' ELSE 'locked' END as status,
			u.created_at,
			u.updated_at,
			u.last_login_at,
			COALESCE(booking_stats.booking_count, 0) as booking_count,
			COALESCE(payment_stats.total_spent, 0) as total_spent
		FROM users u
		LEFT JOIN (
			SELECT user_id, COUNT(*) as booking_count
			FROM bookings
			GROUP BY user_id
		) booking_stats ON u.id = booking_stats.user_id
		LEFT JOIN (
			SELECT payer_id, SUM(amount) as total_spent
			FROM payments
			WHERE status = 'success'
			GROUP BY payer_id
		) payment_stats ON u.id = payment_stats.payer_id
	`

	// Apply filters to the main query
	if req.Status != nil {
		if *req.Status == "locked" {
			userQuery += " WHERE u.deleted_at IS NOT NULL"
		} else {
			userQuery += " WHERE u.deleted_at IS NULL"
		}
	} else {
		userQuery += " WHERE 1=1"
	}

	if req.Role != nil {
		userQuery += fmt.Sprintf(" AND u.role = '%s'", *req.Role)
	}

	if req.Search != nil && *req.Search != "" {
		searchTerm := "%" + *req.Search + "%"
		userQuery += fmt.Sprintf(" AND (u.email LIKE '%s' OR u.full_name LIKE '%s')", searchTerm, searchTerm)
	}

	userQuery += " ORDER BY u.created_at DESC LIMIT ? OFFSET ?"

	err := r.db.Raw(userQuery, limit, offset).Scan(&users).Error
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *adminRepository) LockUser(userID uint64, locked bool, reason *string) error {
	if locked {
		// Soft delete user (lock)
		return r.db.Model(&entity.User{}).Where("id = ?", userID).Update("deleted_at", time.Now()).Error
	} else {
		// Restore user (unlock)
		return r.db.Model(&entity.User{}).Where("id = ?", userID).Update("deleted_at", nil).Error
	}
}

func (r *adminRepository) ApproveVenue(venueID uint64, approved bool, reason *string) error {
	return r.db.Model(&entity.Venue{}).Where("id = ?", venueID).Update("is_active", approved).Error
}

func (r *adminRepository) GetFeedbacksWithDetails(req *entity.AdminFeedbackListRequest) ([]*entity.AdminFeedbackResponse, int64, error) {
	var feedbacks []*entity.AdminFeedbackResponse
	var total int64

	// Build base query
	query := r.db.Model(&entity.Feedback{})

	// Apply filters
	if req.Rating != nil {
		query = query.Where("rating = ?", *req.Rating)
	}

	if req.FromDate != nil && *req.FromDate != "" {
		query = query.Where("DATE(created_at) >= ?", *req.FromDate)
	}

	if req.ToDate != nil && *req.ToDate != "" {
		query = query.Where("DATE(created_at) <= ?", *req.ToDate)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	page := req.Page
	if page <= 0 {
		page = 1
	}
	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}
	offset := (page - 1) * limit

	// Get feedbacks with relationships
	err := query.Preload("FromUser").Preload("ToUser").Preload("Booking").
		Order("created_at DESC").
		Offset(offset).
		Limit(limit).
		Find(&feedbacks).Error

	if err != nil {
		return nil, 0, err
	}

	return feedbacks, total, nil
}
