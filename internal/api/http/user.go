package http

import (
	"saas-template-go/internal/model"
	"saas-template-go/internal/state"
	"saas-template-go/internal/svc"

	"github.com/kataras/iris/v12"
)

func GetUserInfoByID(ctx iris.Context, appState *state.AppState) {
	// user-id from path param.
	userID := ctx.Params().Get("user-id")
	if userID == "" {
		reject(ctx, iris.StatusBadRequest, "缺少 user-id 路径参数")
		return
	}

	// Build service and try to get user info.
	res, err := appState.UserSvc.GetUserInfoByID(userID)
	if err != svc.NO_ERROR {
		reject(ctx, err.Code(), err.Msg())
		return
	}

	accept(ctx, res)
}

func GetUserInfo(ctx iris.Context, appState *state.AppState) {
	// email from query param.
	email := ctx.URLParam("email")
	if email == "" {
		reject(ctx, iris.StatusBadRequest, "缺少 email 参数")
		return
	}

	// Build service and try to get user info.
	res, err := appState.UserSvc.GetUserInfoByEmail(email)
	if err != svc.NO_ERROR {
		reject(ctx, err.Code(), err.Msg())
		return
	}

	accept(ctx, res)
}

func UpdateUserInfo(ctx iris.Context, appState *state.AppState) {
	// user-id from path param.
	userID := ctx.Params().Get("user-id")
	if userID == "" {
		reject(ctx, iris.StatusBadRequest, "缺少 user-id 路径参数")
		return
	}

	// Patch request body.
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
