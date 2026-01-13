package svc

import (
	"fmt"

	"saas-template-go/internal/jwtcodec"
	"saas-template-go/internal/model"
	"saas-template-go/internal/model/po"
	"saas-template-go/internal/repo"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// UserSvc defines service operations for users.
type UserSvc interface {
	GetUserInfoByID(userID string) (SvcRslt[model.UserInfo], SvcErr)
	GetUserInfoByEmail(email string) (SvcRslt[model.UserInfo], SvcErr)
	LoginUser(args model.LoginArgs) (SvcRslt[model.LoginToken], SvcErr)
	UpdateUserInfo(args model.UpdateUserArgs) SvcErr
}

type userSvc struct {
	repo repo.UserRepo
	jwt  *jwtcodec.Codec
}

// NewUserSvc creates a new UserSvc. If r is nil, the default repo implementation is used.
func NewUserSvc(r repo.UserRepo, jwt *jwtcodec.Codec) UserSvc {
	if r == nil {
		panic("UserRepo cannot be nil")
	}

	return &userSvc{repo: r, jwt: jwt}
}

// Get user info by user ID.
func (us *userSvc) GetUserInfoByID(userID string) (SvcRslt[model.UserInfo], SvcErr) {
	userBasics, err := us.repo.GetUserBasicByID(nil, userID)
	if err != nil {
		zap.L().Error("Failed to get user basics by ID", zap.String("userID", userID), zap.Error(err))
		return SvcRslt[model.UserInfo]{}, DB_FAILURE
	}

	userInfo := model.UserInfo{
		UserID:    userBasics.ID,
		Nickname:  userBasics.Nickname,
		CreatedAt: userBasics.CreatedAt,
	}

	return accept(200, userInfo), NO_ERROR
}

// Get user info by email.
func (us *userSvc) GetUserInfoByEmail(email string) (SvcRslt[model.UserInfo], SvcErr) {
	userBasics, err := us.repo.GetUserBasicByEmail(nil, email)
	if err != nil {
		zap.L().Error("Failed to get user basics by email", zap.String("email", email), zap.Error(err))
		return SvcRslt[model.UserInfo]{}, DB_FAILURE
	}

	userInfo := model.UserInfo{
		UserID:    userBasics.ID,
		Nickname:  userBasics.Nickname,
		CreatedAt: userBasics.CreatedAt,
	}

	return accept(200, userInfo), NO_ERROR
}

// Get or create user during login.
func (us *userSvc) LoginUser(args model.LoginArgs) (SvcRslt[model.LoginToken], SvcErr) {
	// Try to get existing user password hash by email.
	// Notice: a optimistic lock based on UNIQUE emial constraint is used here.
	pwdHash, err := us.repo.GetPwdHashByEmail(nil, args.Email)
	if err != nil {
		zap.L().Error("Failed to get user by email during login", zap.String("email", args.Email), zap.Error(err))
		return SvcRslt[model.LoginToken]{}, DB_FAILURE
	}

	if pwdHash == "" {
		// If the user does not exist, register a new user.
		newPwdHash, err := us.hashPwd(args.Password)
		if err != nil {
			zap.L().Error("Failed to hash password during user login", zap.String("email", args.Email), zap.Error(err))
			return SvcRslt[model.LoginToken]{}, PWS_HASH_FAILURE
		}

		newUser := &po.NewUser{
			Email:        args.Email,
			PasswordHash: newPwdHash,
		}

		// The insertion may fail due to UNIQUE email constraint violation.
		// This is expected and acceptable in concurrent login scenarios.
		if err := us.repo.CreateNewUser(nil, newUser); err != nil {
			zap.L().Error("Failed to create new user during login", zap.String("email", args.Email), zap.Error(err))
			return SvcRslt[model.LoginToken]{}, DB_FAILURE
		}

		// Creation succeeded, generate JWT token for the new user.
		// ID is automatically populated after creation.
		token, err := us.generateJWT(newUser.ID)
		if err != nil {
			zap.L().Error("Failed to generate JWT for new user during login", zap.String("email", args.Email), zap.Error(err))
			return SvcRslt[model.LoginToken]{}, DB_FAILURE
		}

		return accept(204, model.LoginToken{Token: token}), NO_ERROR
	}

	// If the user exists, verify the password.
	if !us.verifyPwd(pwdHash, args.Password) {
		return SvcRslt[model.LoginToken]{}, PWD_MISMATCH
	}

	// Verification succeeded, generate JWT token for the existing user.
	token, err := us.generateJWT(args.Email)
	if err != nil {
		zap.L().Error("Failed to generate JWT for existing user during login", zap.String("email", args.Email), zap.Error(err))
		return SvcRslt[model.LoginToken]{}, DB_FAILURE
	}

	return accept(200, model.LoginToken{Token: token}), NO_ERROR
}

// Update user info by user ID.
func (us *userSvc) UpdateUserInfo(args model.UpdateUserArgs) SvcErr {
	// TODO: A better protection logic in updating is needed.
	// Esspecially email field.

	updateUser := &po.UpdateUser{
		ID:       args.UserID,
		Email:    args.Email,
		Nickname: args.Nickname,
	}

	if err := us.repo.UpdateUserByID(nil, updateUser); err != nil {
		zap.L().Error("Failed to update user info", zap.String("userID", args.UserID), zap.Error(err))
		return DB_FAILURE
	}

	return NO_ERROR
}

// Hash a password string using bcrypt with default cost.
func (us *userSvc) hashPwd(pwd string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("Failed to hash password: %w", err)
	}

	return string(hashed), nil
}

// Verify a plain password against a bcrypt hashed password.
func (us *userSvc) verifyPwd(hashedPwd, plainPwd string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPwd)) == nil
}

// Generate a JWT token for a given user ID.
func (us *userSvc) generateJWT(userID string) (string, error) {
	if us.jwt == nil {
		return "", fmt.Errorf("jwt codec is not configured on userSvc")
	}

	token, err := us.jwt.Encode(userID)
	if err != nil {
		return "", fmt.Errorf("Failed to generate JWT: %w", err)
	}

	return token, nil
}
