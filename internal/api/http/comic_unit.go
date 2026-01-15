package http

import (
	"poprako-main-server/internal/model"
	"poprako-main-server/internal/state"
	"poprako-main-server/internal/svc"

	"github.com/kataras/iris/v12"
)

func GetUnitsByPageID(appState *state.AppState) iris.Handler {
	return func(ctx iris.Context) {
		pageID := ctx.Params().Get("page_id")
		if pageID == "" {
			reject(ctx, iris.StatusBadRequest, "缺少 page_id 路径参数")
			return
		}

		res, err := appState.ComicUnitSvc.GetUnitsByPageID(pageID)
		if err != svc.NO_ERROR {
			reject(ctx, err.Code(), err.Msg())
			return
		}

		accept(ctx, res)
	}
}

func CreateUnits(appState *state.AppState) iris.Handler {
	return func(ctx iris.Context) {
		pageID := ctx.Params().Get("page_id")
		if pageID == "" {
			reject(ctx, iris.StatusBadRequest, "缺少 page_id 路径参数")
			return
		}

		var units []model.NewComicUnitArgs

		if err := ctx.ReadJSON(&units); err != nil {
			reject(ctx, iris.StatusBadRequest, "请求体格式错误")
			return
		}

		if len(units) == 0 {
			reject(ctx, iris.StatusBadRequest, "翻译单元列表不能为空")
			return
		}

		// Validate page_id consistency
		for i, u := range units {
			if u.PageID != pageID {
				reject(ctx, iris.StatusBadRequest, "翻译单元的 page_id 与路径参数不一致")
				return
			}
			// Optionally override to ensure consistency
			units[i].PageID = pageID
		}

		opID := ctx.Values().GetString("user_id")
		if opID == "" {
			reject(ctx, iris.StatusUnauthorized, "未认证用户")
			return
		}

		err := appState.ComicUnitSvc.CreateUnits(opID, units)
		if err != svc.NO_ERROR {
			reject(ctx, err.Code(), err.Msg())
			return
		}

		ctx.StatusCode(iris.StatusCreated)
	}
}

func UpdateUnits(appState *state.AppState) iris.Handler {
	return func(ctx iris.Context) {
		var patchUnits []model.PatchComicUnitArgs

		if err := ctx.ReadJSON(&patchUnits); err != nil {
			reject(ctx, iris.StatusBadRequest, "请求体格式错误")
			return
		}

		if len(patchUnits) == 0 {
			reject(ctx, iris.StatusBadRequest, "翻译单元列表不能为空")
			return
		}

		opID := ctx.Values().GetString("user_id")
		if opID == "" {
			reject(ctx, iris.StatusUnauthorized, "未认证用户")
			return
		}

		err := appState.ComicUnitSvc.UpdateUnitsByIDs(opID, patchUnits)
		if err != svc.NO_ERROR {
			reject(ctx, err.Code(), err.Msg())
			return
		}

		ctx.StatusCode(iris.StatusNoContent)
	}
}

func DeleteUnits(appState *state.AppState) iris.Handler {
	return func(ctx iris.Context) {
		var unitIDs []string

		if err := ctx.ReadJSON(&unitIDs); err != nil {
			reject(ctx, iris.StatusBadRequest, "请求体格式错误")
			return
		}

		if len(unitIDs) == 0 {
			reject(ctx, iris.StatusBadRequest, "翻译单元ID列表不能为空")
			return
		}

		err := appState.ComicUnitSvc.DeleteUnitByIDs(unitIDs)
		if err != svc.NO_ERROR {
			reject(ctx, err.Code(), err.Msg())
			return
		}

		ctx.StatusCode(iris.StatusNoContent)
	}
}
