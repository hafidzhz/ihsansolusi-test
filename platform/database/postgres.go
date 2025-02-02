package database

import (
	"fmt"
	"time"

	"github.com/hafidzhz/ihsansolusi-test/app/entity"
	"github.com/hafidzhz/ihsansolusi-test/pkg/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDB() error {
	cfg := config.DBCfg()
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Name,
		cfg.SslMode,
		cfg.Timezone,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	DB.AutoMigrate(
		entity.NewUser(),
	)

	sqlDB, err := DB.DB()

	if err != nil {
		return fmt.Errorf("failed to get generic database object: %w", err)
	}

	sqlDB.SetMaxOpenConns(cfg.MaxOpenConn)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConn)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.MaxConnLifetime) * time.Second)

	return nil
}

func GetDB() *gorm.DB {
	return DB
}
