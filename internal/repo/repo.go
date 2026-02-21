package repo

import (
	"os"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Exct = *gorm.DB

func InitDB() Exct {
	dbURL := os.Getenv("DATABASE_URL")

	if dbURL == "" {
		panic(DB_URL_NOT_SET.Error())
	}

	exec, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		zap.L().Error("Database connection failure", zap.Error(err))
		panic(DB_CONNECTION_FAILURE.Error())
	}

	sqlDB, err := exec.DB()
	if err == nil {
		sqlDB.SetMaxOpenConns(20) // max open
		sqlDB.SetMaxIdleConns(5)  // max idle
	}

	zap.L().Info("Connected to database successfully")

	return exec
}

type Repo interface {
	Exct() Exct
	withTrx(tx Exct) Exct
}

// UnitCounts holds aggregate counts for different unit states on a page.
