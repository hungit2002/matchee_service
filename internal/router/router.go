package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rabbitmq/amqp091-go"
	redislib "github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"matchee/services/internal/config"
	"matchee/services/internal/controller"
	"matchee/services/internal/entity"
	"matchee/services/internal/middleware"
	"matchee/services/internal/repository"
	"matchee/services/internal/service"
	"matchee/services/internal/usecase"
)

func BuildHTTPRouter(cfg *config.Config, logger interface{ Infof(string, ...any) }, db *gorm.DB, rdb *redislib.Client, amqpCh *amqp091.Channel) http.Handler {
	r := gin.New()
	r.Use(gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, entity.OKResponse("Service is healthy", gin.H{
			"app":     cfg.AppName,
			"version": "1.0.0",
			"status":  "running",
		}))
	})

	// Wire repositories
	userRepo := repository.NewUserRepository(db)
	authRepo := repository.NewAuthRepository(db)
	playerRepo := repository.NewPlayerRepository(db)
	venueRepo := repository.NewVenueRepository(db)
	courtRepo := repository.NewCourtRepository(db)
	slotRepo := repository.NewSlotRepository(db)
	matchPostRepo := repository.NewMatchPostRepository(db)
	matchGroupRepo := repository.NewMatchGroupRepository(db)
	matchGroupMemberRepo := repository.NewMatchGroupMemberRepository(db)
	bookingRepo := repository.NewBookingRepository(db)
	paymentRepo := repository.NewPaymentRepository(db)
	notificationRepo := repository.NewNotificationRepository(db)
	feedbackRepo := repository.NewFeedbackRepository(db)
	adminRepo := repository.NewAdminRepository(db)

	// Wire services
	jwtService := service.NewJWTService(cfg.JWTSecret, cfg.JWTExpiry, cfg.RefreshExpiry, authRepo)

	// Wire usecases
	userUC := usecase.NewUserUsecase(userRepo)
	authUC := usecase.NewAuthUsecase(authRepo, userRepo, jwtService)
	playerUC := usecase.NewPlayerUsecase(playerRepo)
	venueUC := usecase.NewVenueUsecase(venueRepo)
	courtUC := usecase.NewCourtUsecase(courtRepo, venueRepo)
	slotUC := usecase.NewSlotUsecase(slotRepo, courtRepo)
	matchPostUC := usecase.NewMatchPostUsecase(matchPostRepo, userRepo, venueRepo)
	matchGroupUC := usecase.NewMatchGroupUsecase(matchGroupRepo, matchGroupMemberRepo, userRepo)
	bookingUC := usecase.NewBookingUsecase(bookingRepo, matchGroupRepo)
	paymentUC := usecase.NewPaymentUsecase(paymentRepo, bookingRepo, userRepo, cfg.PayPalClientID, cfg.PayPalSecret, cfg.PayPalBaseURL)
	notifPublisher := service.NewNotificationPublisher(rdb, amqpCh)
	notificationUC := usecase.NewNotificationUsecase(notificationRepo, notifPublisher)
	feedbackUC := usecase.NewFeedbackUsecase(feedbackRepo, bookingRepo, userRepo)
	adminUC := usecase.NewAdminUsecase(adminRepo, userRepo, venueRepo, feedbackRepo)

	// Wire controllers
	userCtl := controller.NewUserController(userUC)
	authCtl := controller.NewAuthController(authUC)
	playerCtl := controller.NewPlayerController(playerUC)
	venueCtl := controller.NewVenueController(venueUC)
	courtCtl := controller.NewCourtController(courtUC)
	slotCtl := controller.NewSlotController(slotUC)
	matchPostCtl := controller.NewMatchPostController(matchPostUC)
	matchGroupCtl := controller.NewMatchGroupController(matchGroupUC)
	bookingCtl := controller.NewBookingController(bookingUC)
	paymentCtl := controller.NewPaymentController(paymentUC)
	notificationCtl := controller.NewNotificationController(notificationUC)
	feedbackCtl := controller.NewFeedbackController(feedbackUC)
	adminCtl := controller.NewAdminController(adminUC)

	// Wire middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtService)
	adminMiddleware := middleware.NewAdminMiddleware(jwtService)

	// API routes
	api := r.Group("/api")
	v1 := api.Group("/v1")

	// Auth routes (no authentication required)
	auth := v1.Group("/auth")
	auth.POST("/register", authCtl.Register)
	auth.POST("/login", authCtl.Login)
	auth.POST("/refresh", authCtl.RefreshToken)
	auth.POST("/logout", authCtl.Logout)

	// User routes (authentication required)
	users := v1.Group("/users")
	users.Use(authMiddleware.RequireAuth())
	users.GET("/me", authCtl.GetCurrentUser)
	users.PUT("/me", authCtl.UpdateProfile)
	users.POST("/change-password", authCtl.ChangePassword)

	// Player routes (authentication required)
	player := v1.Group("/player")
	player.Use(authMiddleware.RequireAuth())
	player.POST("/profile", playerCtl.CreateOrUpdatePlayerProfile)
	player.GET("/profile", playerCtl.GetPlayerProfileByUserID)
	player.GET("/profile/:id", playerCtl.GetPlayerProfileByID)
	player.GET("/suggestions", playerCtl.GetPlayerSuggestionsQuery)

	// Venue routes
	venues := v1.Group("/venues")
	venues.GET("", venueCtl.GetVenues)
	venues.GET("/:id", venueCtl.GetVenueByID)
	venues.GET("/my", authMiddleware.RequireAuth(), venueCtl.GetMyVenues)
	venues.POST("", authMiddleware.RequireAuth(), venueCtl.CreateVenue)
	venues.PUT("/:id", authMiddleware.RequireAuth(), venueCtl.UpdateVenue)
	venues.DELETE("/:id", authMiddleware.RequireAuth(), venueCtl.DeleteVenue)

	// Court routes
	courts := v1.Group("/courts")
	courts.GET("/:id", courtCtl.GetCourtByID)
	courts.PUT("/:id", authMiddleware.RequireAuth(), courtCtl.UpdateCourt)
	courts.DELETE("/:id", authMiddleware.RequireAuth(), courtCtl.DeleteCourt)

	// Venue-specific court routes
	venues.POST("/:id/courts", authMiddleware.RequireAuth(), courtCtl.CreateCourt)
	venues.GET("/:id/courts", courtCtl.GetCourtsByVenue)

	// Slot routes
	slots := v1.Group("/slots")
	slots.GET("/:id", slotCtl.GetSlotByID)
	slots.PUT("/:id", authMiddleware.RequireAuth(), slotCtl.UpdateSlot)
	slots.DELETE("/:id", authMiddleware.RequireAuth(), slotCtl.DeleteSlot)

	// Court-specific slot routes
	courts.POST("/:id/slots", authMiddleware.RequireAuth(), slotCtl.CreateSlot)
	courts.GET("/:id/slots", slotCtl.GetSlotsByCourt)

	// Match Post routes
	matchPosts := v1.Group("/match-posts")
	matchPosts.POST("", authMiddleware.RequireAuth(), matchPostCtl.CreateMatchPost)
	matchPosts.GET("", matchPostCtl.GetMatchPosts)
	matchPosts.GET("/suggest", matchPostCtl.GetSuggestedMatchPosts)
	matchPosts.GET("/:id", matchPostCtl.GetMatchPostByID)
	matchPosts.PUT("/:id", authMiddleware.RequireAuth(), matchPostCtl.UpdateMatchPost)
	matchPosts.DELETE("/:id", authMiddleware.RequireAuth(), matchPostCtl.DeleteMatchPost)

	// Match Group routes
	matchGroups := v1.Group("/match-groups")
	matchGroups.POST("", matchGroupCtl.CreateMatchGroup)
	matchGroups.GET("", authMiddleware.RequireAuth(), matchGroupCtl.GetMyMatchGroups)
	matchGroups.POST("/:id/members", authMiddleware.RequireAuth(), matchGroupCtl.AddMember)
	matchGroups.DELETE("/:id/members/:user_id", authMiddleware.RequireAuth(), matchGroupCtl.RemoveMember)
	matchGroups.PUT("/:id", authMiddleware.RequireAuth(), matchGroupCtl.UpdateMatchGroup)

	// Booking routes
	bookings := v1.Group("/bookings")
	bookings.POST("", authMiddleware.RequireAuth(), bookingCtl.CreateBooking)
	bookings.GET("", authMiddleware.RequireAuth(), bookingCtl.GetBookings)
	bookings.GET("/:id", authMiddleware.RequireAuth(), bookingCtl.GetBookingByID)
	bookings.PUT("/:id/status", authMiddleware.RequireAuth(), bookingCtl.UpdateBookingStatus)
	bookings.DELETE("/:id", authMiddleware.RequireAuth(), bookingCtl.DeleteBooking)

	// Payment routes
	payments := v1.Group("/payments")
	payments.POST("", authMiddleware.RequireAuth(), paymentCtl.CreatePayment)
	payments.GET("", authMiddleware.RequireAuth(), paymentCtl.GetPayments)
	payments.GET("/:id", authMiddleware.RequireAuth(), paymentCtl.GetPaymentByID)
	payments.POST("/webhook", paymentCtl.HandlePayPalWebhook) // No auth required for webhook

	// Notification routes
	notif := v1.Group("/notifications")
	notif.GET("", authMiddleware.RequireAuth(), notificationCtl.GetMyNotifications)
	notif.PUT("/:id/read", authMiddleware.RequireAuth(), notificationCtl.MarkRead)
	notif.POST("", authMiddleware.RequireAuth(), notificationCtl.Create)

	// Feedback routes
	feedbacks := v1.Group("/feedbacks")
	feedbacks.POST("", authMiddleware.RequireAuth(), feedbackCtl.CreateFeedback)
	feedbacks.GET("/user/:id", feedbackCtl.GetUserFeedbacks)
	feedbacks.GET("/booking/:id", feedbackCtl.GetBookingFeedbacks)

	// Admin routes
	admin := v1.Group("/admin")
	admin.Use(adminMiddleware.RequireAdmin())
	admin.GET("/statistics/revenue", adminCtl.GetRevenueStatistics)
	admin.GET("/users", adminCtl.GetUsers)
	admin.PUT("/users/:id/lock", adminCtl.LockUser)
	admin.PUT("/venues/:id/approve", adminCtl.ApproveVenue)
	admin.GET("/feedbacks", adminCtl.GetFeedbacks)

	// Legacy user routes (for backward compatibility)
	legacyUsers := api.Group("/users")
	legacyUsers.POST("", userCtl.Register)
	legacyUsers.GET(":id", userCtl.Get)

	return r
}
