package http

import (
	"poprako-main-server/internal/model"
	"poprako-main-server/internal/state"
	"poprako-main-server/internal/svc"

	"github.com/kataras/iris/v12"
)

func GetPageByID(appState *state.AppState) iris.Handler {
	return func(ctx iris.Context) {
		pageID := ctx.Params().Get("page_id")
		if pageID == "" {
			reject(ctx, iris.StatusBadRequest, "缺少 page_id 路径参数")
			return
		}

		res, err := appState.ComicPageSvc.GetPageByID(pageID)
		if err != svc.NO_ERROR {
			reject(ctx, err.Code(), err.Msg())
			return
		}

		accept(ctx, res)
	}
}

func GetPagesByComicID(appState *state.AppState) iris.Handler {
	return func(ctx iris.Context) {
		comicID := ctx.Params().Get("comic_id")
		if comicID == "" {
			reject(ctx, iris.StatusBadRequest, "缺少 comic_id 路径参数")
			return
		}

		res, err := appState.ComicPageSvc.GetPagesByComicID(comicID)
		if err != svc.NO_ERROR {
			reject(ctx, err.Code(), err.Msg())
			return
		}

		accept(ctx, res)
	}
}

func CreatePages(appState *state.AppState) iris.Handler {
	return func(ctx iris.Context) {
		var pages []model.CreateComicPageArgs

		if err := ctx.ReadJSON(&pages); err != nil {
			reject(ctx, iris.StatusBadRequest, "请求体格式错误")
			return
		}

		if len(pages) == 0 {
			reject(ctx, iris.StatusBadRequest, "页面列表不能为空")
			return
		}

		opID := ctx.Values().GetString("user_id")
		if opID == "" {
			reject(ctx, iris.StatusUnauthorized, "未认证用户")
			return
		}

		res, err := appState.ComicPageSvc.CreatePages(opID, pages)
		if err != svc.NO_ERROR {
			reject(ctx, err.Code(), err.Msg())
			return
		}

		accept(ctx, res)
	}
}

func UpdatePageByID(appState *state.AppState) iris.Handler {
	return func(ctx iris.Context) {
		pageID := ctx.Params().Get("page_id")
		if pageID == "" {
			reject(ctx, iris.StatusBadRequest, "缺少 page_id 路径参数")
			return
		}

		var args model.PatchComicPageArgs

		if err := ctx.ReadJSON(&args); err != nil {
			reject(ctx, iris.StatusBadRequest, "请求体格式错误")
			return
		}

		args.ID = pageID

		opID := ctx.Values().GetString("user_id")
		if opID == "" {
			reject(ctx, iris.StatusUnauthorized, "未认证用户")
			return
		}

		err := appState.ComicPageSvc.UpdatePageByID(opID, &args)
		if err != svc.NO_ERROR {
			reject(ctx, err.Code(), err.Msg())
			return
		}

		ctx.StatusCode(iris.StatusNoContent)
	}
}

func DeletePageByID(appState *state.AppState) iris.Handler {
	return func(ctx iris.Context) {
		pageID := ctx.Params().Get("page_id")
		if pageID == "" {
			reject(ctx, iris.StatusBadRequest, "缺少 page_id 路径参数")
			return
		}

		err := appState.ComicPageSvc.DeletePageByID(pageID)
		if err != svc.NO_ERROR {
			reject(ctx, err.Code(), err.Msg())
			return
		}

		ctx.StatusCode(iris.StatusNoContent)
	}
}
