package http

import (
	"poprako-main-server/internal/model"
	"poprako-main-server/internal/state"
	"poprako-main-server/internal/svc"

	"github.com/kataras/iris/v12"
)

func GetComicInfoByID(appState *state.AppState) iris.Handler {
	return func(ctx iris.Context) {
		comicID := ctx.Params().Get("comic_id")
		if comicID == "" {
			reject(ctx, iris.StatusBadRequest, "缺少 comic_id 路径参数")
			return
		}

		res, err := appState.ComicSvc.GetComicInfoByID(comicID)
		if err != svc.NO_ERROR {
			reject(ctx, err.Code(), err.Msg())
			return
		}

		accept(ctx, res)
	}
}

func GetComicBriefsByWorksetID(appState *state.AppState) iris.Handler {
	return func(ctx iris.Context) {
		worksetID := ctx.Params().Get("workset_id")
		if worksetID == "" {
			reject(ctx, iris.StatusBadRequest, "缺少 workset_id 路径参数")
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

		res, err := appState.ComicSvc.GetComicBriefsByWorksetID(worksetID, opt.Offset, opt.Limit)
		if err != svc.NO_ERROR {
			reject(ctx, err.Code(), err.Msg())
			return
		}

		accept(ctx, res)
	}
}

func RetrieveComicBriefs(appState *state.AppState) iris.Handler {
	return func(ctx iris.Context) {
		var opt model.RetrieveComicOpt

		if err := ctx.ReadQuery(&opt); err != nil {
			reject(ctx, iris.StatusBadRequest, "查询参数格式错误")
			return
		}

		res, err := appState.ComicSvc.RetrieveComics(opt)
		if err != svc.NO_ERROR {
			reject(ctx, err.Code(), err.Msg())
			return
		}

		accept(ctx, res)
	}
}

func CreateComic(appState *state.AppState) iris.Handler {
	return func(ctx iris.Context) {
		var args model.CreateComicArgs

		if err := ctx.ReadJSON(&args); err != nil {
			reject(ctx, iris.StatusBadRequest, "请求体格式错误")
			return
		}

		opID := ctx.Values().GetString("user_id")
		if opID == "" {
			reject(ctx, iris.StatusUnauthorized, "未认证用户")
			return
		}

		res, err := appState.ComicSvc.CreateComic(opID, args)
		if err != svc.NO_ERROR {
			reject(ctx, err.Code(), err.Msg())
			return
		}

		accept(ctx, res)
	}
}

func UpdateComicByID(appState *state.AppState) iris.Handler {
	return func(ctx iris.Context) {
		comicID := ctx.Params().Get("comic_id")
		if comicID == "" {
			reject(ctx, iris.StatusBadRequest, "缺少 comic_id 路径参数")
			return
		}

		var args model.UpdateComicArgs

		if err := ctx.ReadJSON(&args); err != nil {
			reject(ctx, iris.StatusBadRequest, "请求体格式错误")
			return
		}

		args.ID = comicID

		err := appState.ComicSvc.UpdateComicByID(args)
		if err != svc.NO_ERROR {
			reject(ctx, err.Code(), err.Msg())
			return
		}

		ctx.StatusCode(iris.StatusNoContent)
	}
}

func DeleteComicByID(appState *state.AppState) iris.Handler {
	return func(ctx iris.Context) {
		comicID := ctx.Params().Get("comic_id")
		if comicID == "" {
			reject(ctx, iris.StatusBadRequest, "缺少 comic_id 路径参数")
			return
		}

		err := appState.ComicSvc.DeleteComicByID(comicID)
		if err != svc.NO_ERROR {
			reject(ctx, err.Code(), err.Msg())
			return
		}

		ctx.StatusCode(iris.StatusNoContent)
	}
}
