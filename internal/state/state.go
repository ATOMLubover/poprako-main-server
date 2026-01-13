// Package state provides application-wide state management.
// Objects and variables that need be to injected throughout the application
// are to be stored here.
package state

import (
	"saas-template-go/internal/config"
	"saas-template-go/internal/jwtcodec"
	"saas-template-go/internal/svc"
)

type AppState struct {
	Cfg      config.AppCfg
	JWTCodec *jwtcodec.Codec
	UserSvc  svc.UserSvc
}

func NewAppState(cfg config.AppCfg, jwtCodec *jwtcodec.Codec, userSvc svc.UserSvc) AppState {
	return AppState{
		Cfg:      cfg,
		JWTCodec: jwtCodec,
		UserSvc:  userSvc,
	}
}
