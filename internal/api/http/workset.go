package http

import (
	"poprako-main-server/internal/model"
	"poprako-main-server/internal/state"
	"poprako-main-server/internal/svc"

	"github.com/kataras/iris/v12"
)

func GetWorksetByID(appState *state.AppState) iris.Handler {
	return func(ctx iris.Context) {
		worksetID := ctx.Params().Get("workset_id")
		if worksetID == "" {
			reject(ctx, iris.StatusBadRequest, "缺少 workset_id 路径参数")
			return
		}

		res, err := appState.WorksetSvc.GetWorksetByID(worksetID)
		if err != svc.NO_ERROR {
			reject(ctx, err.Code(), err.Msg())
			return
		}

		accept(ctx, res)
	}
}

func RetrieveWorksets(appState *state.AppState) iris.Handler {
	return func(ctx iris.Context) {
		var opt struct {
			Limit  int `url:"limit,default=10"`
			Offset int `url:"offset,default=0"`
		}

		if err := ctx.ReadQuery(&opt); err != nil {
			reject(ctx, iris.StatusBadRequest, "查询参数格式错误")
			return
		}

		res, err := appState.WorksetSvc.RetrieveWorksets(opt.Limit, opt.Offset)
		if err != svc.NO_ERROR {
			reject(ctx, err.Code(), err.Msg())
			return
		}

		accept(ctx, res)
	}
}

func CreateWorkset(appState *state.AppState) iris.Handler {
	return func(ctx iris.Context) {
		var args model.CreateWorksetArgs

		if err := ctx.ReadJSON(&args); err != nil {
			reject(ctx, iris.StatusBadRequest, "请求体格式错误")
			return
		}

		opID := ctx.Values().GetString("user_id")
		if opID == "" {
			reject(ctx, iris.StatusUnauthorized, "未认证用户")
			return
		}

		res, err := appState.WorksetSvc.CreateWorkset(opID, &args)
		if err != svc.NO_ERROR {
			reject(ctx, err.Code(), err.Msg())
			return
		}

		accept(ctx, res)
	}
}

func UpdateWorksetByID(appState *state.AppState) iris.Handler {
	return func(ctx iris.Context) {
		worksetID := ctx.Params().Get("workset_id")
		if worksetID == "" {
			reject(ctx, iris.StatusBadRequest, "缺少 workset_id 路径参数")
			return
		}

		var args model.UpdateWorksetArgs

		if err := ctx.ReadJSON(&args); err != nil {
			reject(ctx, iris.StatusBadRequest, "请求体格式错误")
			return
		}

		args.ID = worksetID

		err := appState.WorksetSvc.UpdateWorksetByID(&args)
		if err != svc.NO_ERROR {
			reject(ctx, err.Code(), err.Msg())
			return
		}

		ctx.StatusCode(iris.StatusNoContent)
	}
}
