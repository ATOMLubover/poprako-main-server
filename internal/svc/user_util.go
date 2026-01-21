package svc

import (
	"errors"
	"fmt"
	"strconv"

	"poprako-main-server/internal/model"
	"poprako-main-server/internal/repo"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// Role strings
const (
	ROLE_TRANSLATOR  = "translator"
	ROLE_PROOFREADER = "proofreader"
	ROLE_TYPESETTER  = "typesetter"
	ROLE_REDRAWER    = "redrawer"
	ROLE_REVIEWER    = "reviewer"
	ROLE_UPLOADER    = "uploader"
	ROLE_ADMIN       = "admin"
)

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
func (us *userSvc) verifyInvCode(code string, qq string) (*model.InvitationInfo, error) {
	invitation, err := us.invRepo.GetInvitationByQQ(nil, qq)
	if err == repo.REC_NOT_FOUND {
		return nil, errors.New("invitation record not found")
	}
	if err != nil {
		zap.L().Error("Failed to get invitation by qq during invitation code verification", zap.String("qq", qq), zap.Error(err))
		return nil, err
	}

	if invitation.InvCode != code {
		return nil, errors.New("invitation code does not match")
	}

	return &model.InvitationInfo{
		InvitorID:         invitation.InvitorID,
		InviteeQQ:         invitation.InviteeQQ,
		Pending:           invitation.Pending,
		AssignTranslator:  invitation.AssignTranslator,
		AssignProofreader: invitation.AssignProofreader,
		AssignTypesetter:  invitation.AssignTypesetter,
		AssignRedrawer:    invitation.AssignRedrawer,
		AssignReviewer:    invitation.AssignReviewer,
		AssignUploader:    invitation.AssignUploader,
	}, nil
}
