// Package main provides the entry point for the server application.
package main

import (
	"fmt"

	"saas-template-go/internal/api/http"
	"saas-template-go/internal/config"
	"saas-template-go/internal/jwtcodec"
	"saas-template-go/internal/logger"
	"saas-template-go/internal/repo"
	"saas-template-go/internal/state"
	"saas-template-go/internal/svc"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	// Load environment variables from .env file.
	// Logger and JWT codec will depend on these variables.
	initEnv()

	initLogger()

	cfg := config.LoadConfig("")

	state := initAppState(cfg)

	http.Run(state)
}

func initEnv() {
	if err := godotenv.Load(); err != nil {
		panic(fmt.Sprintf("Error loading .env file: %v", err))
	}
}

func initLogger() {
	// Initialize global logger.
	lgr := logger.InitLogger()

	zap.ReplaceGlobals(lgr)
}

func initAppState(cfg config.AppCfg) state.AppState {
	// Create JWT codec.
	jwtCodec := jwtcodec.NewJWTCodec(cfg.JWTExpSecs)

	// Create repositories.
	ex := repo.InitDB()

	userRepo := repo.NewUserRepo(ex)

	// Create services.
	userSvc := svc.NewUserSvc(userRepo, jwtCodec)

	return state.NewAppState(cfg, jwtCodec, userSvc)
}
