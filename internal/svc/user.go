package svc

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"poprako-main-server/internal/jwtcodec"
	"poprako-main-server/internal/model"
	"poprako-main-server/internal/model/po"
	"poprako-main-server/internal/repo"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// UserSvc defines service operations for users.
type UserSvc interface {
	GetUserInfoByID(userID string) (SvcRslt[model.UserInfo], SvcErr)
	GetUserInfoByQQ(qq string) (SvcRslt[model.UserInfo], SvcErr)

	LoginUser(args model.LoginArgs) (SvcRslt[model.LoginReply], SvcErr)

	UpdateUserInfo(args model.UpdateUserArgs) SvcErr
}

type userSvc struct {
	repo repo.UserRepo
	jwt  *jwtcodec.Codec

	mu       sync.Mutex
	invCodes map[string]struct{}
}

// NewUserSvc creates a new UserSvc. If r is nil, the default repo implementation is used.
func NewUserSvc(r repo.UserRepo, jwt *jwtcodec.Codec) UserSvc {
	if r == nil {
		panic("UserRepo cannot be nil")
	}

	return &userSvc{
		repo:     r,
		jwt:      jwt,
		invCodes: make(map[string]struct{}),
	}
}

// Get user info by user ID.
func (us *userSvc) GetUserInfoByID(userID string) (SvcRslt[model.UserInfo], SvcErr) {
	userBasic, err := us.repo.GetUserByID(nil, userID)
	if err != nil {
		zap.L().Error("Failed to get user basics by ID", zap.String("userID", userID), zap.Error(err))
		return SvcRslt[model.UserInfo]{}, DB_FAILURE
	}

	userInfo := model.UserInfo{
		UserID:    userBasic.ID,
		IsAdmin:   userBasic.IsAdmin,
		Nickname:  userBasic.Nickname,
		CreatedAt: userBasic.CreatedAt,
	}
	if userBasic.AssignedTranslatorAt != nil {
		userInfo.AssignedTranslatorAt = *userBasic.AssignedTranslatorAt
	}
	if userBasic.AssignedProofreaderAt != nil {
		userInfo.AssignedProofreaderAt = *userBasic.AssignedProofreaderAt
	}
	if userBasic.AssignedTypesetterAt != nil {
		userInfo.AssignedTypesetterAt = *userBasic.AssignedTypesetterAt
	}
	if userBasic.AssignedRedrawerAt != nil {
		userInfo.AssignedRedrawerAt = *userBasic.AssignedRedrawerAt
	}
	if userBasic.AssignedReviewerAt != nil {
		userInfo.AssignedReviewerAt = *userBasic.AssignedReviewerAt
	}
	if userBasic.AssignedUploaderAt != nil {
		userInfo.AssignedUploaderAt = *userBasic.AssignedUploaderAt
	}

	return accept(200, userInfo), NO_ERROR
}

// Get user info by QQ.
func (us *userSvc) GetUserInfoByQQ(qq string) (SvcRslt[model.UserInfo], SvcErr) {
	userBasics, err := us.repo.GetUserByQQ(nil, qq)
	if err != nil {
		zap.L().Error("Failed to get user basics by email", zap.String("email", qq), zap.Error(err))
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
func (us *userSvc) LoginUser(args model.LoginArgs) (SvcRslt[model.LoginReply], SvcErr) {
	if args.InvCode != "" {
		// If invitation code is provided, validate it.
		// This happens when a new user is trying to register.

		invCode, err := us.verifyInvCode(args.InvCode)
		if err != nil {
			zap.L().Warn("Invalid invitation code during user login", zap.String("invCode", args.InvCode), zap.Error(err))
			return SvcRslt[model.LoginReply]{}, INV_CODE_INVALID
		}

		if invCode != args.QQ {
			zap.L().Error("Invitation is not corresponding with qq", zap.String("invCode", invCode), zap.String("qq", args.QQ))
			return SvcRslt[model.LoginReply]{}, INV_CODE_MISMATCH
		}

		// Verification passed.

		newPwdHash, err := us.hashPwd(args.Password)
		if err != nil {
			zap.L().Error("Failed to hash password during user login", zap.String("qq", args.QQ), zap.Error(err))
			return SvcRslt[model.LoginReply]{}, PWD_HASH_FAILURE
		}

		newID, err := genUUID()
		if err != nil {
			zap.L().Error("Failed to generate UUID during user login", zap.String("qq", args.QQ), zap.Error(err))
			return SvcRslt[model.LoginReply]{}, ID_GEN_FAILURE
		}

		newUser := &po.NewUser{
			ID:           newID,
			QQ:           args.QQ,
			Nickname:     args.Nickname,
			PasswordHash: newPwdHash,
		}

		// The insertion may fail due to UNIQUE email constraint violation.
		// This is expected and acceptable in concurrent login scenarios.
		if err := us.repo.CreateUser(nil, newUser); err != nil {
			zap.L().Error("Failed to create new user during login", zap.String("qq", args.QQ), zap.Error(err))
			return SvcRslt[model.LoginReply]{}, DB_FAILURE
		}

		// Creation succeeded, generate JWT token for the new user.
		// ID is automatically populated after creation.
		token, err := us.genJWT(newUser.ID)
		if err != nil {
			zap.L().Error("Failed to generate JWT for new user during login", zap.String("qq", args.QQ), zap.Error(err))
			return SvcRslt[model.LoginReply]{}, DB_FAILURE
		}

		return accept(204, model.LoginReply{Token: token}), NO_ERROR
	}

	// No invitation code provided, treat as existing user login.

	secret, err := us.repo.GetSecretUserByQQ(nil, args.QQ)
	if err != nil {
		zap.L().Error("Failed to get password hash during user login", zap.String("email", args.QQ), zap.Error(err))
		return SvcRslt[model.LoginReply]{}, DB_FAILURE
	}

	// If the user exists, verify the password.
	if !us.verifyPwd(secret.PwdHash, args.Password) {
		return SvcRslt[model.LoginReply]{}, PWD_MISMATCH
	}

	// Verification succeeded, generate JWT token for the existing user.
	token, err := us.genJWT(secret.ID)
	if err != nil {
		zap.L().Error("Failed to generate JWT for existing user during login", zap.String("qq", args.QQ), zap.Error(err))
		return SvcRslt[model.LoginReply]{}, DB_FAILURE
	}

	return accept(200, model.LoginReply{Token: token}), NO_ERROR
}

// Update user info by user ID.
func (us *userSvc) UpdateUserInfo(args model.UpdateUserArgs) SvcErr {
	updateUser := &po.PatchUser{
		ID:       args.UserID,
		QQ:       args.QQ,
		Nickname: args.Nickname,
	}

	// Handle assignment fields.
	var zero int64 = 0

	if args.AssignTranslator != nil {
		if *args.AssignTranslator {
			now := time.Now().Unix()
			updateUser.AssignedTranslatorAt = &now
		} else {
			updateUser.AssignedTranslatorAt = &zero
		}
	}
	if args.AssignProofreader != nil {
		if *args.AssignProofreader {
			now := time.Now().Unix()
			updateUser.AssignedProofreaderAt = &now
		} else {
			updateUser.AssignedProofreaderAt = &zero
		}
	}
	if args.AssignTypesetter != nil {
		if *args.AssignTypesetter {
			now := time.Now().Unix()
			updateUser.AssignedTypesetterAt = &now
		} else {
			updateUser.AssignedTypesetterAt = &zero
		}
	}
	if args.AssignRedrawer != nil {
		if *args.AssignRedrawer {
			now := time.Now().Unix()
			updateUser.AssignedRedrawerAt = &now
		} else {
			updateUser.AssignedRedrawerAt = &zero
		}
	}
	if args.AssignReviewer != nil {
		if *args.AssignReviewer {
			now := time.Now().Unix()
			updateUser.AssignedReviewerAt = &now
		} else {
			updateUser.AssignedReviewerAt = &zero
		}
	}
	if args.AssignUploader != nil {
		if *args.AssignUploader {
			now := time.Now().Unix()
			updateUser.AssignedUploaderAt = &now
		} else {
			updateUser.AssignedUploaderAt = &zero
		}
	}

	if err := us.repo.UpdateUserByID(nil, updateUser); err != nil {
		zap.L().Error("Failed to update user info", zap.String("userID", args.UserID), zap.Error(err))
		return DB_FAILURE
	}

	return NO_ERROR
}

func (us *userSvc) InviteUser(operUserID string, args model.InviteUserArgs) (SvcRslt[model.InviteUserReply], SvcErr) {
	userBasic, err := us.repo.GetUserByID(nil, operUserID)
	if err != nil {
		zap.L().Error("Failed to get user basics by ID during invitation", zap.String("userID", operUserID), zap.Error(err))
		return SvcRslt[model.InviteUserReply]{}, DB_FAILURE
	}

	// Check whether operator is admin.
	if !userBasic.IsAdmin {
		zap.L().Warn("Non-admin user attempted to invite user", zap.String("userID", operUserID))
		return SvcRslt[model.InviteUserReply]{}, PERMISSION_DENIED
	}

	// Generate invitation code.
	invCode := us.genInvCode(args.InviteeID)
	if invCode == "" {
		zap.L().Error("Failed to generate invitation code", zap.String("inviteeID", args.InviteeID))
		return SvcRslt[model.InviteUserReply]{}, DB_FAILURE
	}

	return accept(200, model.InviteUserReply{InvCode: invCode}), NO_ERROR
}

// Utility functions.

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
func (us *userSvc) genJWT(userID string) (string, error) {
	if us.jwt == nil {
		zap.L().Warn("us.jwt not set in generateJWT")
		return "", fmt.Errorf("jwt codec is not configured on userSvc")
	}

	token, err := us.jwt.Encode(userID)
	if err != nil {
		return "", fmt.Errorf("Failed to generate JWT: %w", err)
	}

	return token, nil
}

func (us *userSvc) genInvCode(decStr string) string {
	// Encode invitation code: from dec string to hex string.
	num, err := strconv.ParseInt(decStr, 10, 32)
	if err != nil {
		zap.L().Error("Failed to parse invitation code number", zap.String("decStr", decStr), zap.Error(err))
		return ""
	}

	hexStr := strconv.FormatInt(num, 16)

	us.mu.Lock()
	defer us.mu.Unlock()

	if len(us.invCodes) >= 50 {
		// Limit the number of stored invitation codes up to 50.
		zap.L().Warn("Failed add any more invitation code due to capacity issues")
		return ""
	}

	us.invCodes[hexStr] = struct{}{}

	return hexStr
}

// Check whether a invitation code is valid.
func (us *userSvc) verifyInvCode(codeStr string) (string, error) {
	us.mu.Lock()

	if _, exists := us.invCodes[codeStr]; !exists {
		// Invitation code does not exist.
		us.mu.Unlock()

		return "", errors.New("invitation code invalid")
	}

	// Mark the invitation code as used.
	delete(us.invCodes, codeStr)

	us.mu.Unlock()

	// Decode invitation code: from hex string to dec string.
	hexNum, err := strconv.ParseInt(codeStr, 16, 32)
	if err != nil {
		return "", errors.New("Failed to parse invitation code")
	}

	decStr := strconv.FormatInt(hexNum, 10)

	return decStr, nil
}
