package http

import (
	"fmt"

	"poprako-main-server/internal/model"
	"poprako-main-server/internal/state"

	"github.com/kataras/iris/v12"
)

type versionTuple struct {
	Major int
	Minor int
	Patch int
}

func (v *versionTuple) LessThan(other versionTuple) bool {
	if v.Major != other.Major {
		return v.Major < other.Major
	}

	if v.Minor != other.Minor {
		return v.Minor < other.Minor
	}

	return v.Patch < other.Patch
}

func parseVersion(versionStr string) (versionTuple, error) {
	var v versionTuple

	n, err := fmt.Sscanf(versionStr, "%d.%d.%d", &v.Major, &v.Minor, &v.Patch)
	if err != nil || n != 3 {
		return versionTuple{}, fmt.Errorf("invalid version format: %s", versionStr)
	}

	return v, nil
}

func marshalVersion(v versionTuple) string {
	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

func CheckUpdateHandler(appState *state.AppState) iris.Handler {
	return func(ctx iris.Context) {
		cliAppVersion := ctx.GetHeader("X-Client-App-Version")

		cliAppVersionTuple, err := parseVersion(cliAppVersion)
		if err != nil {
			reject(ctx, iris.StatusBadRequest, "无效的客户端版本格式")
			return
		}

		minAppVersionTuple, err := parseVersion(appState.Cfg.CheckUpdateMin)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			return
		}

		allowUsage := !cliAppVersionTuple.LessThan(minAppVersionTuple)

		res := model.CheckVersionReply{
			LatestVersion: appState.Cfg.NativeAppVersion,
			Title:         appState.Cfg.CheckUpdateTitle,
			Description:   appState.Cfg.CheckUpdateDesc,
			AllowUsage:    allowUsage,
		}

		ctx.JSON(res)
		ctx.StatusCode(iris.StatusOK)
	}
}
