package svc

// Unified wrapper for service results.
type SvcRslt[T any] struct {
	Code uint16
	Msg  string
	Data *T
}

// Convenience functions for creating ServResult instances
// when returning a failed response.
func accept[T any](code uint16, data T) SvcRslt[T] {
	return SvcRslt[T]{
		Code: code,
		Data: &data,
	}
}

// Convenience functions for creating ServResult instances
// when returning a successful response.
func reject[T any](code uint16, msg string) SvcRslt[T] {
	return SvcRslt[T]{
		Code: code,
		Msg:  msg,
	}
}

// Define error enums for service results.
type SvcErr string

// Detailed service error enums.
const (
	// Placeholder for no error.
	NO_ERROR SvcErr = ""
	// Generic database failure.
	DB_FAILURE SvcErr = "Database failure"
	// Password hashing failure.
	PWS_HASH_FAILURE SvcErr = "Password hashing failure"
	// User not found.
	USER_NOT_FOUND SvcErr = "User not found"
	// Password mismatch.
	PWD_MISMATCH SvcErr = "Password mismatch"
	// User ID mismatch in args and path.
	USER_ID_MISMATCH SvcErr = "User ID mismatch"
)

// Get a API error code for the ServError.
func (e *SvcErr) Code() uint16 {
	switch *e {
	case NO_ERROR:
		return 200
	case DB_FAILURE:
		return 500
	case PWS_HASH_FAILURE:
		return 500
	case USER_NOT_FOUND:
		return 404
	case PWD_MISMATCH:
		return 401
	default:
		return 500
	}
}

// Get the error message in Chinese for the ServError.
// NOTICE: it shadows the specific error message in the enum,
// so it is only suitable for API response messages.
func (e *SvcErr) Msg() string {
	switch *e {
	case NO_ERROR:
		return ""
	case DB_FAILURE:
		return "服务器内部错误"
	case PWS_HASH_FAILURE:
		return "服务器内部错误"
	case USER_NOT_FOUND:
		return "用户不存在"
	case PWD_MISMATCH:
		return "密码错误"
	default:
		return "服务器内部错误"
	}
}
