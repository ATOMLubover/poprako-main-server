// Package main provides the entry point for the server application.
package main

import (
	"fmt"

	"poprako-main-server/internal/api/http"
	"poprako-main-server/internal/config"
	"poprako-main-server/internal/jwtcodec"
	"poprako-main-server/internal/logger"
	"poprako-main-server/internal/oss"
	"poprako-main-server/internal/repo"
	"poprako-main-server/internal/state"
	"poprako-main-server/internal/svc"

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
	comicRepo := repo.NewComicRepo(ex)
	worksetRepo := repo.NewWorksetRepo(ex)
	comicUnitRepo := repo.NewComicUnitRepo(ex)
	comicAsgnRepo := repo.NewComicAsgnRepo(ex)
	comicPageRepo := repo.NewComicPageRepo(ex)

	// Create OSS client.
	ossClient := oss.NewR2Client()

	// Create services.
	userSvc := svc.NewUserSvc(userRepo, jwtCodec)
	comicSvc := svc.NewComicSvc(comicRepo, userRepo, comicAsgnRepo)
	worksetSvc := svc.NewWorksetSvc(worksetRepo, userRepo)
	comicUnitSvc := svc.NewComicUnitSvc(comicUnitRepo)
	comicAsgnSvc := svc.NewComicAsgnSvc(comicAsgnRepo)
	comicPageSvc := svc.NewComicPageSvc(comicPageRepo, comicRepo, ossClient)

	return state.NewAppState(
		cfg,
		jwtCodec,
		userSvc,
		comicSvc,
		worksetSvc,
		comicUnitSvc,
		comicAsgnSvc,
		comicPageSvc,
		ossClient,
	)
}
