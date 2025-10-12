package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rabbitmq/amqp091-go"
	redislib "github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"matchee/services/internal/config"
	"matchee/services/internal/controller"
	"matchee/services/internal/repository"
	"matchee/services/internal/usecase"
)

func BuildHTTPRouter(cfg *config.Config, logger interface{ Infof(string, ...any) }, db *gorm.DB, rdb *redislib.Client, amqpCh *amqp091.Channel) http.Handler {
	r := gin.New()
	r.Use(gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "app": cfg.AppName})
	})

	// Wire users
	userRepo := repository.NewUserRepository(db)
	userUC := usecase.NewUserUsecase(userRepo)
	userCtl := controller.NewUserController(userUC)

	api := r.Group("/api")
	users := api.Group("/users")
	users.POST("", userCtl.Register)
	users.GET(":id", userCtl.Get)

	return r
}
