package svc

import (
	"sync"
	"time"

	"poprako-main-server/internal/jwtcodec"
	"poprako-main-server/internal/model"
	"poprako-main-server/internal/model/po"
	"poprako-main-server/internal/repo"

	"go.uber.org/zap"
)

// UserSvc defines service operations for users.
type UserSvc interface {
	GetUserInfoByID(userID string) (SvcRslt[model.UserInfo], SvcErr)
	GetUserInfoByQQ(qq string) (SvcRslt[model.UserInfo], SvcErr)

	InviteUser(operUserID string, args model.InviteUserArgs) (SvcRslt[model.InviteUserReply], SvcErr)
	LoginUser(args model.LoginArgs) (SvcRslt[model.LoginReply], SvcErr)

	UpdateUserInfo(args model.UpdateUserArgs) SvcErr
	AssignUserRole(opID string, args model.AssignUserRoleArgs) SvcErr

	GetUserInfos(opt model.RetrieveUserOpt) (SvcRslt[[]model.UserInfo], SvcErr)
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
		CreatedAt: userBasic.CreatedAt.Unix(),
	}
	if userBasic.AssignedTranslatorAt != nil {
		ts := userBasic.AssignedTranslatorAt.Unix()
		userInfo.AssignedTranslatorAt = &ts
	}
	if userBasic.AssignedProofreaderAt != nil {
		ts := userBasic.AssignedProofreaderAt.Unix()
		userInfo.AssignedProofreaderAt = &ts
	}
	if userBasic.AssignedTypesetterAt != nil {
		ts := userBasic.AssignedTypesetterAt.Unix()
		userInfo.AssignedTypesetterAt = &ts
	}
	if userBasic.AssignedRedrawerAt != nil {
		ts := userBasic.AssignedRedrawerAt.Unix()
		userInfo.AssignedRedrawerAt = &ts
	}
	if userBasic.AssignedReviewerAt != nil {
		ts := userBasic.AssignedReviewerAt.Unix()
		userInfo.AssignedReviewerAt = &ts
	}
	if userBasic.AssignedUploaderAt != nil {
		ts := userBasic.AssignedUploaderAt.Unix()
		userInfo.AssignedUploaderAt = &ts
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
		CreatedAt: userBasics.CreatedAt.Unix(),
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
		zap.L().Warn("Failed to get password hash during user login", zap.String("email", args.QQ), zap.Error(err))
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

	// // Handle assignment fields.
	// zero := int64(0)

	// if args.AssignTranslator != nil {
	// 	if *args.AssignTranslator {
	// 		now := time.Now().Unix()
	// 		updateUser.AssignedTranslatorAt = &now
	// 	} else {
	// 		updateUser.AssignedTranslatorAt = &zero
	// 	}
	// }
	// if args.AssignProofreader != nil {
	// 	if *args.AssignProofreader {
	// 		now := time.Now().Unix()
	// 		updateUser.AssignedProofreaderAt = &now
	// 	} else {
	// 		updateUser.AssignedProofreaderAt = &zero
	// 	}
	// }
	// if args.AssignTypesetter != nil {
	// 	if *args.AssignTypesetter {
	// 		now := time.Now().Unix()
	// 		updateUser.AssignedTypesetterAt = &now
	// 	} else {
	// 		updateUser.AssignedTypesetterAt = &zero
	// 	}
	// }
	// if args.AssignRedrawer != nil {
	// 	if *args.AssignRedrawer {
	// 		now := time.Now().Unix()
	// 		updateUser.AssignedRedrawerAt = &now
	// 	} else {
	// 		updateUser.AssignedRedrawerAt = &zero
	// 	}
	// }
	// if args.AssignReviewer != nil {
	// 	if *args.AssignReviewer {
	// 		now := time.Now().Unix()
	// 		updateUser.AssignedReviewerAt = &now
	// 	} else {
	// 		updateUser.AssignedReviewerAt = &zero
	// 	}
	// }
	// if args.AssignUploader != nil {
	// 	if *args.AssignUploader {
	// 		now := time.Now().Unix()
	// 		updateUser.AssignedUploaderAt = &now
	// 	} else {
	// 		updateUser.AssignedUploaderAt = &zero
	// 	}
	// }

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

func (us *userSvc) GetUserInfos(opt model.RetrieveUserOpt) (SvcRslt[[]model.UserInfo], SvcErr) {
	userBasics, err := us.repo.RetrieveUsers(nil, opt)
	if err != nil {
		zap.L().Error("Failed to retrieve user basics", zap.Error(err))
		return SvcRslt[[]model.UserInfo]{}, DB_FAILURE
	}

	userInfos := make([]model.UserInfo, 0, len(userBasics))

	for _, ub := range userBasics {
		ui := model.UserInfo{
			UserID:    ub.ID,
			QQ:        ub.QQ,
			IsAdmin:   ub.IsAdmin,
			Nickname:  ub.Nickname,
			CreatedAt: ub.CreatedAt.Unix(),
		}
		if ub.AssignedTranslatorAt != nil {
			ts := ub.AssignedTranslatorAt.Unix()
			ui.AssignedTranslatorAt = &ts
		}
		if ub.AssignedProofreaderAt != nil {
			ts := ub.AssignedProofreaderAt.Unix()
			ui.AssignedProofreaderAt = &ts
		}
		if ub.AssignedTypesetterAt != nil {
			ts := ub.AssignedTypesetterAt.Unix()
			ui.AssignedTypesetterAt = &ts
		}
		if ub.AssignedRedrawerAt != nil {
			ts := ub.AssignedRedrawerAt.Unix()
			ui.AssignedRedrawerAt = &ts
		}
		if ub.AssignedReviewerAt != nil {
			ts := ub.AssignedReviewerAt.Unix()
			ui.AssignedReviewerAt = &ts
		}
		if ub.AssignedUploaderAt != nil {
			ts := ub.AssignedUploaderAt.Unix()
			ui.AssignedUploaderAt = &ts
		}

		userInfos = append(userInfos, ui)
	}

	return accept(200, userInfos), NO_ERROR
}

func (us *userSvc) AssignUserRole(
	opID string,
	args model.AssignUserRoleArgs,
) SvcErr {
	var patch po.PatchUser

	hasChange := false

	now := time.Now()
	zero := time.Time{}

	roles := args.Roles

	for _, r := range roles {
		switch r.Role {
		case ROLE_TRANSLATOR:
			hasChange = true

			if r.Assigned {
				patch.AssignedTranslatorAt = &now
			} else {
				patch.AssignedTranslatorAt = &zero
			}

		case ROLE_PROOFREADER:
			hasChange = true

			if r.Assigned {
				patch.AssignedProofreaderAt = &now
			} else {
				patch.AssignedProofreaderAt = &zero
			}

		case ROLE_TYPESETTER:
			hasChange = true

			if r.Assigned {
				patch.AssignedTypesetterAt = &now
			} else {
				patch.AssignedTypesetterAt = &zero
			}

		case ROLE_REDRAWER:
			hasChange = true

			if r.Assigned {
				patch.AssignedRedrawerAt = &now
			} else {
				patch.AssignedRedrawerAt = &zero
			}

		case ROLE_REVIEWER:
			hasChange = true

			if r.Assigned {
				patch.AssignedReviewerAt = &now
			} else {
				patch.AssignedReviewerAt = &zero
			}

		case ROLE_UPLOADER:
			hasChange = true

			if r.Assigned {
				patch.AssignedUploaderAt = &now
			} else {
				patch.AssignedUploaderAt = &zero
			}

		default:
			zap.L().Warn("Unknown role string in AssignUserRole", zap.String("role", r.Role))
			return INVALID_ROLE_DATA
		}
	}

	if !hasChange {
		zap.L().Warn("No valid role changes in AssignUserRole", zap.String("userID", args.UserID), zap.Any("roles", roles))
		return INVALID_ROLE_DATA
	}

	patch.ID = args.UserID

	if err := us.repo.UpdateUserByID(nil, &patch); err != nil {
		zap.L().Error("Failed to update user roles in AssignUserRole", zap.String("userID", args.UserID), zap.Error(err))
		return DB_FAILURE
	}

	return NO_ERROR
}
