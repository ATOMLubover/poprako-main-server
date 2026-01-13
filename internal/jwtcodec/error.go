package jwtcodec

// Define error enums for JWT codec operations.
type JWTError = string

const (
	SECRET_NOT_SET            JWTError = "Secret key not set"
	SIGN_FAILURE              JWTError = "Failed to sign token"
	UNEXPECTED_SIGNING_METHOD JWTError = "Unexpected signing method"
	INVALID_TOKEN             JWTError = "Invalid token"
	CLAIMS_PARSING_FAILURE    JWTError = "Invalid claims"
)
