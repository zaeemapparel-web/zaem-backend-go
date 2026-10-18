package database

import (
	"time"

	"github.com/rs/zerolog/log"
	"github.com/zaem/backend-go/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() {
	var err error

	logLevel := gormlogger.Warn
	if config.C.IsDevelopment() {
		logLevel = gormlogger.Warn // Info for SQL logs
	}

	DB, err = gorm.Open(postgres.Open(config.C.DatabaseURL), &gorm.Config{
		Logger: gormlogger.Default.LogMode(logLevel),
		NowFunc: func() time.Time { return time.Now().UTC() },
		SkipDefaultTransaction:                   true,
		DisableForeignKeyConstraintWhenMigrating: true,
		PrepareStmt:                              true,
	})

	if err != nil {
		log.Fatal().Err(err).Msg("❌ Failed to connect to PostgreSQL")
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal().Err(err).Msg("❌ Failed to get sql.DB")
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxLifetime(time.Hour)
	sqlDB.SetConnMaxIdleTime(10 * time.Minute)

	if err := sqlDB.Ping(); err != nil {
		log.Fatal().Err(err).Msg("❌ Database ping failed")
	}

	log.Info().Msg("✅ Database connected (PostgreSQL via GORM)")
}

func Close() {
	if DB == nil {
		return
	}
	sqlDB, err := DB.DB()
	if err != nil {
		return
	}
	_ = sqlDB.Close()
	log.Info().Msg("🔌 Database connection closed")
}

func Health() error {
	if DB == nil {
		return nil
	}
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Ping()
}