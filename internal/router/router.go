package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rabbitmq/amqp091-go"
	redislib "github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"matchee/services/internal/config"
	"matchee/services/internal/controller"
	"matchee/services/internal/middleware"
	"matchee/services/internal/repository"
	"matchee/services/internal/service"
	"matchee/services/internal/usecase"
)

func BuildHTTPRouter(cfg *config.Config, logger interface{ Infof(string, ...any) }, db *gorm.DB, rdb *redislib.Client, amqpCh *amqp091.Channel) http.Handler {
	r := gin.New()
	r.Use(gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "app": cfg.AppName})
	})

	// Wire repositories
	userRepo := repository.NewUserRepository(db)
	authRepo := repository.NewAuthRepository(db)
	playerRepo := repository.NewPlayerRepository(db)

	// Wire services
	jwtService := service.NewJWTService(cfg.JWTSecret, cfg.JWTExpiry, cfg.RefreshExpiry, authRepo)

	// Wire usecases
	userUC := usecase.NewUserUsecase(userRepo)
	authUC := usecase.NewAuthUsecase(authRepo, userRepo, jwtService)
	playerUC := usecase.NewPlayerUsecase(playerRepo)

	// Wire controllers
	userCtl := controller.NewUserController(userUC)
	authCtl := controller.NewAuthController(authUC)
	playerCtl := controller.NewPlayerController(playerUC)

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

	// Legacy user routes (for backward compatibility)
	legacyUsers := api.Group("/users")
	legacyUsers.POST("", userCtl.Register)
	legacyUsers.GET(":id", userCtl.Get)

	return r
}
