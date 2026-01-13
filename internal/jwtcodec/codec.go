package jwtcodec

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

// UserClaims contains the claims stored in JWT tokens.
type UserClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// Codec handles encoding and decoding of JWT tokens.
// At current stage, it uses a simple SHA-256 HMAC with a secret key.
type Codec struct {
	expSecs   int64
	secretKey string
}

func NewJWTCodec(expSecs int64) *Codec {
	sec := os.Getenv("JWT_SECRET_KEY")

	if sec == "" {
		panic(SECRET_NOT_SET)
	}

	return &Codec{
		expSecs:   expSecs,
		secretKey: sec,
	}
}

func (jc *Codec) Encode(userID string) (string, error) {
	now := time.Now()
	exp := now.Add(time.Duration(jc.expSecs) * time.Second)

	claims := &UserClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
			// TODO: Add more registered claims if needed, e.g., Issuer, Subject, Audience.
		},
	}

	// Use HS256 signing method to generate the token.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string.
	tokenStr, err := token.SignedString([]byte(jc.secretKey))
	if err != nil {
		zap.L().Error("Failed to sign user claims", zap.Error(err))
		return "", errors.New(SIGN_FAILURE)
	}

	return tokenStr, nil
}

func (jc *Codec) Decode(tokenStr string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&UserClaims{},
		func(token *jwt.Token) (any, error) {
			// Validate the signing method. HS256 is expected.
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				zap.L().Warn("Unexpected signing method", zap.Any("alg", token.Header["alg"]))
				return nil, errors.New(UNEXPECTED_SIGNING_METHOD)
			}

			return []byte(jc.secretKey), nil
		},
	)
	if err != nil {
		zap.L().Error("Failed to parse token", zap.Error(err))
		return nil, errors.New(INVALID_TOKEN)
	}

	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New(CLAIMS_PARSING_FAILURE)
}
