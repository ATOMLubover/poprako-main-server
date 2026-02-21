package svc

import (
	"poprako-main-server/internal/model"
	"poprako-main-server/internal/model/po"
	"poprako-main-server/internal/repo"

	"go.uber.org/zap"
)

type InvitationSvc interface {
	GetInvitationInfos(opID string) (SvcRslt[[]model.InvitationInfo], SvcErr)

	CreateInvitation(opID string, args model.CreateInvitationArgs) (SvcRslt[model.CreateInvitationReply], SvcErr)
}

type invitationSvc struct {
	userRepo repo.UserRepo
	invRepo  repo.InvitationRepo
}

func NewInvitationSvc(invRepo repo.InvitationRepo, userRepo repo.UserRepo) InvitationSvc {
	return &invitationSvc{
		invRepo:  invRepo,
		userRepo: userRepo,
	}
}

func (is *invitationSvc) GetInvitationInfos(
	opID string,
) (
	SvcRslt[[]model.InvitationInfo],
	SvcErr,
) {
	opBasic, err := is.userRepo.GetUserByID(nil, opID)
	if err != nil {
		zap.L().Error("Failed to get user basics by ID during invitation", zap.String("userID", opID), zap.Error(err))
		return SvcRslt[[]model.InvitationInfo]{}, DB_FAILURE
	}

	// Check whether operator is admin.
	if !opBasic.IsAdmin {
		zap.L().Warn("Non-admin user attempted to invite user", zap.String("userID", opID))
		return SvcRslt[[]model.InvitationInfo]{}, PERMISSION_DENIED
	}

	invitations, err := is.invRepo.RetrieveInvitations(nil)
	if err != nil {
		return SvcRslt[[]model.InvitationInfo]{}, DB_FAILURE
	}

	var infos []model.InvitationInfo

	for _, inv := range invitations {
		infos = append(infos, model.InvitationInfo{
			ID:                inv.ID,
			InvitorID:         inv.InvitorID,
			InviteeQQ:         inv.InviteeQQ,
			InvCode:           inv.InvCode,
			AssignTranslator:  inv.AssignTranslator,
			AssignProofreader: inv.AssignProofreader,
			AssignTypesetter:  inv.AssignTypesetter,
			AssignRedrawer:    inv.AssignRedrawer,
			AssignReviewer:    inv.AssignReviewer,
			AssignUploader:    inv.AssignUploader,
			Pending:           inv.Pending,
			CreatedAt:         inv.CreatedAt.Unix(),
		})
	}

	return accept(200, infos), NO_ERROR
}

func (is *invitationSvc) CreateInvitation(
	opID string,
	args model.CreateInvitationArgs,
) (
	SvcRslt[model.CreateInvitationReply],
	SvcErr,
) {
	opBasic, err := is.userRepo.GetUserByID(nil, opID)
	if err != nil {
		zap.L().Error("Failed to get user basics by ID during invitation", zap.String("userID", opID), zap.Error(err))
		return SvcRslt[model.CreateInvitationReply]{}, DB_FAILURE
	}

	// Check whether operator is admin.
	if !opBasic.IsAdmin {
		zap.L().Warn("Non-admin user attempted to invite user", zap.String("userID", opID))
		return SvcRslt[model.CreateInvitationReply]{}, PERMISSION_DENIED
	}

	// Verification passed.

	// Check whether invitee QQ is already registered.
	_, err = is.userRepo.GetUserByQQ(nil, args.InviteeQQ)
	if err != nil && err != repo.REC_NOT_FOUND {
		zap.L().Error("Failed to check existing user by QQ during invitation", zap.String("qq", args.InviteeQQ), zap.Error(err))
		return SvcRslt[model.CreateInvitationReply]{}, DB_FAILURE
	}

	newID, err := genUUID()
	if err != nil {
		zap.L().Error("Failed to generate UUID for new invitation", zap.Error(err))
		return SvcRslt[model.CreateInvitationReply]{}, ID_GEN_FAILURE
	}

	invCode := genInvCode(newID)

	newInvitation := &po.NewInvitation{
		ID:        newID,
		InvitorID: opID,
		InviteeQQ: args.InviteeQQ,
		InvCode:   invCode,
	}

	if args.AssignTranslator != nil {
		newInvitation.AssignTranslator = *args.AssignTranslator
	}
	if args.AssignProofreader != nil {
		newInvitation.AssignProofreader = *args.AssignProofreader
	}
	if args.AssignTypesetter != nil {
		newInvitation.AssignTypesetter = *args.AssignTypesetter
	}
	if args.AssignRedrawer != nil {
		newInvitation.AssignRedrawer = *args.AssignRedrawer
	}
	if args.AssignReviewer != nil {
		newInvitation.AssignReviewer = *args.AssignReviewer
	}
	if args.AssignUploader != nil {
		newInvitation.AssignUploader = *args.AssignUploader
	}

	if err := is.invRepo.CreateInvitations(nil, newInvitation); err != nil {
		return SvcRslt[model.CreateInvitationReply]{}, DB_FAILURE
	}

	return accept(200, model.CreateInvitationReply{
		InvCode: invCode,
	}), NO_ERROR
}

func genInvCode(newID string) string {
	// To simplify, just use UUID
	return newID[len(newID)-6:]
}
