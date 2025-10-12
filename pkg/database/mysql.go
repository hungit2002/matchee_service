package database

import (
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"matchee/services/internal/config"
)

func NewMySQL(cfg *config.Config) (*gorm.DB, error) {
	dialector := mysql.Open(cfg.MySQLDSN)
	return gorm.Open(dialector, &gorm.Config{
		Logger:  logger.Default.LogMode(logger.Warn),
		NowFunc: func() time.Time { return time.Now().UTC() },
	})
}
