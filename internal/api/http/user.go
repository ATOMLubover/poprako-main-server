package http

import (
	"fmt"
	"net/http"

	"poprako-main-server/internal/model"
	"poprako-main-server/internal/state"
	"poprako-main-server/internal/svc"

	"github.com/kataras/iris/v12"
)

func GetCurrUserInfo(appState *state.AppState) iris.Handler {
	return func(ctx iris.Context) {
		userID := ctx.Values().GetString("user_id")
		if userID == "" {
			reject(ctx, iris.StatusUnauthorized, "未认证用户")
			return
		}

		res, err := appState.UserSvc.GetUserInfoByID(userID)
		if err != svc.NO_ERROR {
			reject(ctx, err.Code(), err.Msg())
			return
		}

		accept(ctx, res)
	}
}

func GetUserInfoByID(appState *state.AppState) iris.Handler {
	return func(ctx iris.Context) {
		// user_id from path param.
		userID := ctx.Params().Get("user_id")
		if userID == "" {
			reject(ctx, iris.StatusBadRequest, "缺少 user_id 路径参数")
			return
		}

		res, err := appState.UserSvc.GetUserInfoByID(userID)
		if err != svc.NO_ERROR {
			reject(ctx, err.Code(), err.Msg())
			return
		}

		accept(ctx, res)
	}
}

func UpdateUserInfo(appState *state.AppState) iris.Handler {
	return func(ctx iris.Context) {
		// user_id from path param.
		userID := ctx.Params().Get("user_id")
		if userID == "" {
			reject(ctx, iris.StatusBadRequest, "缺少 user_id 路径参数")
			return
		}

		var args model.UpdateUserArgs

		if err := ctx.ReadJSON(&args); err != nil {
			reject(ctx, iris.StatusBadRequest, "请求体格式错误")
			return
		}

		// Check user-id consistency.
		if userID != args.UserID {
			reject(ctx, iris.StatusBadRequest, "路径参数 user-id 与请求体内 user_id 不匹配")
			return
		}

		// Build service and try to update user info.
		err := appState.UserSvc.UpdateUserInfo(args)
		if err != svc.NO_ERROR {
			reject(ctx, err.Code(), err.Msg())
			return
		}

		// Return 204 if successful.
		ctx.StatusCode(iris.StatusNoContent)
	}
}

func InviteUser(appState *state.AppState) iris.Handler {
	return func(ctx iris.Context) {
		var args model.CreateInvitationArgs

		if err := ctx.ReadJSON(&args); err != nil {
			reject(ctx, iris.StatusBadRequest, "请求体格式错误")
			return
		}

		opID := ctx.Values().GetString("user_id")
		if opID == "" {
			reject(ctx, iris.StatusUnauthorized, "未认证用户")
			return
		}

		res, err := appState.InvitationSvc.CreateInvitation(opID, args)
		if err != svc.NO_ERROR {
			reject(ctx, err.Code(), err.Msg())
			return
		}

		accept(ctx, res)
	}
}

func GetInvitations(appState *state.AppState) iris.Handler {
	return func(ctx iris.Context) {
		opID := ctx.Values().GetString("user_id")
		if opID == "" {
			reject(ctx, iris.StatusUnauthorized, "未认证用户")
			return
		}

		res, err := appState.InvitationSvc.GetInvitationInfos(opID)
		if err != svc.NO_ERROR {
			reject(ctx, err.Code(), err.Msg())
			return
		}

		accept(ctx, res)
	}
}

func LoginUser(appState *state.AppState) iris.Handler {
	return func(ctx iris.Context) {
		var args model.LoginArgs

		if err := ctx.ReadJSON(&args); err != nil {
			reject(ctx, iris.StatusBadRequest, "请求体格式错误")
			return
		}

		res, err := appState.UserSvc.LoginUser(args)
		if err != svc.NO_ERROR {
			reject(ctx, err.Code(), err.Msg())
			return
		}

		ctx.SetCookie(&http.Cookie{
			Name:  "Authorization",
			Value: fmt.Sprintf("Bearer %s", res.Data.Token),
		})

		accept(ctx, res)
	}
}

func RetrieveUserInfos(appState *state.AppState) iris.Handler {
	return func(ctx iris.Context) {
		var opt model.RetrieveUserOpt

		if err := ctx.ReadQuery(&opt); err != nil {
			reject(ctx, iris.StatusBadRequest, "查询参数格式错误")
			return
		}

		res, err := appState.UserSvc.GetUserInfos(opt)
		if err != svc.NO_ERROR {
			reject(ctx, err.Code(), err.Msg())
			return
		}

		accept(ctx, res)
	}
}

func AssignUserRole(appState *state.AppState) iris.Handler {
	return func(ctx iris.Context) {
		opID := ctx.Values().GetString("user_id")
		if opID == "" {
			reject(ctx, iris.StatusUnauthorized, "未认证用户")
			return
		}

		userID := ctx.Params().Get("user_id")
		if userID == "" {
			reject(ctx, iris.StatusBadRequest, "缺少 user_id 路径参数")
			return
		}

		var args model.AssignUserRoleArgs

		if err := ctx.ReadJSON(&args); err != nil {
			reject(ctx, iris.StatusBadRequest, "请求体格式错误")
			return
		}

		if userID != args.UserID {
			reject(ctx, iris.StatusBadRequest, "路径参数 user_id 与请求体内 user_id 不匹配")
			return
		}

		err := appState.UserSvc.AssignUserRole(opID, args)
		if err != svc.NO_ERROR {
			reject(ctx, err.Code(), err.Msg())
			return
		}

		ctx.StatusCode(iris.StatusNoContent)
	}
}
