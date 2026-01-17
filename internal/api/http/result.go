package http

import (
	"poprako-main-server/internal/svc"

	"github.com/kataras/iris/v12"
)

// HTTPRslt is a unified wrapper for HTTP responses.
// When it represents a successful response, Data is populated and Msg is empty.
// When it represents a failed response, Msg is populated and Data is nil.
type HTTPRslt[T any] struct {
	Code uint16 `json:"code"`
	Msg  string `json:"msg,omitempty"`
	Data *T     `json:"data,omitempty"`
}

// Unified wrapper to response to failed request.
func reject(ctx iris.Context, code uint16, msg string) {
	ctx.StatusCode(int(code))

	ctx.JSON(HTTPRslt[string]{
		Code: code,
		Msg:  msg,
	})
}

// A quick transform function that converts
// SvcRslt into HTTPResult and writes it to the context.
func accept[T any](ctx iris.Context, res svc.SvcRslt[T]) {
	ctx.JSON(HTTPRslt[T]{
		Code: res.Code,
		Msg:  res.Msg,
		Data: res.Data,
	})
}
