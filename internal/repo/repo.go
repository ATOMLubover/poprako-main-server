package repo

import (
	"os"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Executor = *gorm.DB

func InitDB() Executor {
	dbURL := os.Getenv("DATABASE_URL")

	if dbURL == "" {
		panic(DB_URL_NOT_SET.Error())
	}

	exec, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})
	if err != nil {
		zap.L().Error("Database connection failure", zap.Error(err))
		panic(DB_CONNECTION_FAILURE.Error())
	}

	zap.L().Info("Connected to database successfully")

	return exec
}

type Repo interface {
	Exec() Executor
	withTrx(tx Executor) Executor
}

// UnitCounts holds aggregate counts for different unit states on a page.
type UnitCounts struct {
	Inbox      int64
	Outbox     int64
	Translated int64
	Proved     int64
}
