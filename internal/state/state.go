// Package state provides application-wide state management.
// Objects and variables that need be to injected throughout the application
// are to be stored here.
package state

import (
	"poprako-main-server/internal/config"
	"poprako-main-server/internal/jwtcodec"
	"poprako-main-server/internal/oss"
	"poprako-main-server/internal/svc"
)

type AppState struct {
	Cfg          config.AppCfg
	JWTCodec     *jwtcodec.Codec
	UserSvc      svc.UserSvc
	ComicSvc     svc.ComicSvc
	WorksetSvc   svc.WorksetSvc
	ComicUnitSvc svc.ComicUnitSvc
	ComicAsgnSvc svc.ComicAsgnSvc
	ComicPageSvc svc.ComicPageSvc
	OSSClient    oss.OSSClient
}

func NewAppState(
	cfg config.AppCfg,
	jwtCodec *jwtcodec.Codec,
	userSvc svc.UserSvc,
	comicSvc svc.ComicSvc,
	worksetSvc svc.WorksetSvc,
	comicUnitSvc svc.ComicUnitSvc,
	comicAsgnSvc svc.ComicAsgnSvc,
	comicPageSvc svc.ComicPageSvc,
	ossClient oss.OSSClient,
) AppState {
	return AppState{
		Cfg:          cfg,
		JWTCodec:     jwtCodec,
		UserSvc:      userSvc,
		ComicSvc:     comicSvc,
		WorksetSvc:   worksetSvc,
		ComicUnitSvc: comicUnitSvc,
		ComicAsgnSvc: comicAsgnSvc,
		ComicPageSvc: comicPageSvc,
		OSSClient:    ossClient,
	}
}
