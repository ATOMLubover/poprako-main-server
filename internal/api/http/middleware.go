package http

import (
	"strings"

	"poprako-main-server/internal/state"

	"github.com/kataras/iris/v12"
)

// AuthMiddleware validates JWT token and sets user_id in context
func AuthMiddleware(appState *state.AppState) iris.Handler {
	return func(ctx iris.Context) {
		// Get Authorization header
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			reject(ctx, iris.StatusUnauthorized, "缺少认证令牌")
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			reject(ctx, iris.StatusUnauthorized, "认证令牌格式错误")
			return
		}

		token := parts[1]

		// Decode and validate token
		claims, err := appState.JWTCodec.Decode(token)
		if err != nil {
			reject(ctx, iris.StatusUnauthorized, "认证令牌无效或已过期")
			return
		}

		// Set user_id in context for downstream handlers
		ctx.Values().Set("user_id", claims.UserID)

		// Continue to next handler
		ctx.Next()
	}
}
