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

	// Wire services
	jwtService := service.NewJWTService(cfg.JWTSecret, cfg.JWTExpiry, cfg.RefreshExpiry, authRepo)

	// Wire usecases
	userUC := usecase.NewUserUsecase(userRepo)
	authUC := usecase.NewAuthUsecase(authRepo, userRepo, jwtService)
	playerUC := usecase.NewPlayerUsecase(playerRepo)
	venueUC := usecase.NewVenueUsecase(venueRepo)
	courtUC := usecase.NewCourtUsecase(courtRepo, venueRepo)
	slotUC := usecase.NewSlotUsecase(slotRepo, courtRepo)

	// Wire controllers
	userCtl := controller.NewUserController(userUC)
	authCtl := controller.NewAuthController(authUC)
	playerCtl := controller.NewPlayerController(playerUC)
	venueCtl := controller.NewVenueController(venueUC)
	courtCtl := controller.NewCourtController(courtUC)
	slotCtl := controller.NewSlotController(slotUC)

	// Wire middleware
	authMiddleware := middleware.NewAuthMiddleware(jwtService)

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
	venues.GET("", venueCtl.GetVenues)                                        // Public: Get venues with filters
	venues.GET("/:id", venueCtl.GetVenueByID)                                 // Public: Get venue details
	venues.GET("/my", authMiddleware.RequireAuth(), venueCtl.GetMyVenues)     // Private: Get my venues
	venues.POST("", authMiddleware.RequireAuth(), venueCtl.CreateVenue)       // Private: Create venue
	venues.PUT("/:id", authMiddleware.RequireAuth(), venueCtl.UpdateVenue)    // Private: Update venue
	venues.DELETE("/:id", authMiddleware.RequireAuth(), venueCtl.DeleteVenue) // Private: Delete venue

	// Court routes
	courts := v1.Group("/courts")
	courts.GET("/:id", courtCtl.GetCourtByID)                                 // Public: Get court details
	courts.PUT("/:id", authMiddleware.RequireAuth(), courtCtl.UpdateCourt)    // Private: Update court
	courts.DELETE("/:id", authMiddleware.RequireAuth(), courtCtl.DeleteCourt) // Private: Delete court

	// Venue-specific court routes
	venues.POST("/:id/courts", authMiddleware.RequireAuth(), courtCtl.CreateCourt) // Private: Create court for venue
	venues.GET("/:id/courts", courtCtl.GetCourtsByVenue)                           // Public: Get courts for venue

	// Slot routes
	slots := v1.Group("/slots")
	slots.GET("/:id", slotCtl.GetSlotByID)                                 // Public: Get slot details
	slots.PUT("/:id", authMiddleware.RequireAuth(), slotCtl.UpdateSlot)    // Private: Update slot
	slots.DELETE("/:id", authMiddleware.RequireAuth(), slotCtl.DeleteSlot) // Private: Delete slot

	// Court-specific slot routes
	courts.POST("/:id/slots", authMiddleware.RequireAuth(), slotCtl.CreateSlot) // Private: Create slot for court
	courts.GET("/:id/slots", slotCtl.GetSlotsByCourt)                           // Public: Get slots for court

	// Legacy user routes (for backward compatibility)
	legacyUsers := api.Group("/users")
	legacyUsers.POST("", userCtl.Register)
	legacyUsers.GET(":id", userCtl.Get)

	return r
}
