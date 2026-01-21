package http

import (
	"poprako-main-server/internal/model"
	"poprako-main-server/internal/state"
	"poprako-main-server/internal/svc"

	"github.com/kataras/iris/v12"
)

func GetAsgnByID(appState *state.AppState) iris.Handler {
	return func(ctx iris.Context) {
		asgnID := ctx.Params().Get("asgn_id")
		if asgnID == "" {
			reject(ctx, iris.StatusBadRequest, "缺少 asgn_id 路径参数")
			return
		}

		res, err := appState.ComicAsgnSvc.GetAsgnByID(asgnID)
		if err != svc.NO_ERROR {
			reject(ctx, err.Code(), err.Msg())
			return
		}

		accept(ctx, res)
	}
}

func GetAsgnsByComicID(appState *state.AppState) iris.Handler {
	return func(ctx iris.Context) {
		comicID := ctx.Params().Get("comic_id")
		if comicID == "" {
			reject(ctx, iris.StatusBadRequest, "缺少 comic_id 路径参数")
			return
		}

		var opt struct {
			Limit  int `url:"limit,default=10"`
			Offset int `url:"offset,default=0"`
		}

		if err := ctx.ReadQuery(&opt); err != nil {
			reject(ctx, iris.StatusBadRequest, "查询参数格式错误")
			return
		}

		res, err := appState.ComicAsgnSvc.GetAsgnsByComicID(comicID, opt.Offset, opt.Limit)
		if err != svc.NO_ERROR {
			reject(ctx, err.Code(), err.Msg())
			return
		}

		accept(ctx, res)
	}
}

func GetAsgnsByUserID(appState *state.AppState) iris.Handler {
	return func(ctx iris.Context) {
		userID := ctx.Params().Get("user_id")
		if userID == "" {
			reject(ctx, iris.StatusBadRequest, "缺少 user_id 路径参数")
			return
		}

		var opt struct {
			Limit  int `url:"limit,default=10"`
			Offset int `url:"offset,default=0"`
		}

		if err := ctx.ReadQuery(&opt); err != nil {
			reject(ctx, iris.StatusBadRequest, "查询参数格式错误")
			return
		}

		res, err := appState.ComicAsgnSvc.GetAsgnsByUserID(userID, opt.Offset, opt.Limit)
		if err != svc.NO_ERROR {
			reject(ctx, err.Code(), err.Msg())
			return
		}

		accept(ctx, res)
	}
}

func CreateAsgn(appState *state.AppState) iris.Handler {
	return func(ctx iris.Context) {
		var args model.CreateComicAsgnArgs

		if err := ctx.ReadJSON(&args); err != nil {
			reject(ctx, iris.StatusBadRequest, "请求体格式错误")
			return
		}

		res, err := appState.ComicAsgnSvc.CreateAsgn(args)
		if err != svc.NO_ERROR {
			reject(ctx, err.Code(), err.Msg())
			return
		}

		accept(ctx, res)
	}
}

func UpdateAsgn(appState *state.AppState) iris.Handler {
	return func(ctx iris.Context) {
		asgnID := ctx.Params().Get("asgn_id")
		if asgnID == "" {
			reject(ctx, iris.StatusBadRequest, "缺少 asgn_id 路径参数")
			return
		}

		var args model.UpdateComicAsgnArgs

		if err := ctx.ReadJSON(&args); err != nil {
			reject(ctx, iris.StatusBadRequest, "请求体格式错误")
			return
		}

		args.ID = asgnID

		err := appState.ComicAsgnSvc.UpdateAsgnByID(args)
		if err != svc.NO_ERROR {
			reject(ctx, err.Code(), err.Msg())
			return
		}

		ctx.StatusCode(iris.StatusNoContent)
	}
}

func DeleteAsgnByID(appState *state.AppState) iris.Handler {
	return func(ctx iris.Context) {
		asgnID := ctx.Params().Get("asgn_id")
		if asgnID == "" {
			reject(ctx, iris.StatusBadRequest, "缺少 asgn_id 路径参数")
			return
		}

		err := appState.ComicAsgnSvc.DeleteAsgnByID(asgnID)
		if err != svc.NO_ERROR {
			reject(ctx, err.Code(), err.Msg())
			return
		}

		ctx.StatusCode(iris.StatusNoContent)
	}
}
