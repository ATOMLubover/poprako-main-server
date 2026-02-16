package svc

import (
	"poprako-main-server/internal/model"
	"poprako-main-server/internal/model/po"
	"poprako-main-server/internal/repo"

	"go.uber.org/zap"
)

// WorksetSvc defines service operations for worksets.
type WorksetSvc interface {
	GetWorksetByID(worksetID string) (SvcRslt[model.WorksetInfo], SvcErr)
	RetrieveWorksets(limit, offset int) (SvcRslt[[]model.WorksetInfo], SvcErr)

	CreateWorkset(opID string, args *model.CreateWorksetArgs) (SvcRslt[model.CreateWorksetReply], SvcErr)

	UpdateWorksetByID(args *model.UpdateWorksetArgs) SvcErr

	DeleteWorksetByID(worksetID string) SvcErr
}

type worksetSvc struct {
	repo     repo.WorksetRepo
	userRepo repo.UserRepo
}

// NewWorksetSvc creates a new WorksetSvc. r must not be nil.
func NewWorksetSvc(r repo.WorksetRepo, ur repo.UserRepo) WorksetSvc {
	if r == nil {
		panic("WorksetRepo cannot be nil")
	}
	if ur == nil {
		panic("UserRepo cannot be nil")
	}

	return &worksetSvc{repo: r, userRepo: ur}
}

func (ws *worksetSvc) GetWorksetByID(worksetID string) (SvcRslt[model.WorksetInfo], SvcErr) {
	detail, err := ws.repo.GetWorksetByID(nil, worksetID)
	if err != nil {
		zap.L().Error("Failed to get workset by ID", zap.String("worksetID", worksetID), zap.Error(err))
		return SvcRslt[model.WorksetInfo]{}, DB_FAILURE
	}

	info := model.WorksetInfo{
		ID:              detail.ID,
		Index:           detail.Index,
		Name:            detail.Name,
		ComicCount:      detail.ComicCount,
		Description:     detail.Description,
		CreatorID:       detail.CreatorID,
		CreatorNickname: detail.CreatorNickname,
		CreatedAt:       detail.CreatedAt.Unix(),
		UpdatedAt:       detail.UpdatedAt.Unix(),
	}

	return accept(200, info), NO_ERROR
}

func (ws *worksetSvc) RetrieveWorksets(limit, offset int) (SvcRslt[[]model.WorksetInfo], SvcErr) {
	details, err := ws.repo.RetrieveWorksets(nil, limit, offset)
	if err != nil {
		zap.L().Error("Failed to retrieve worksets", zap.Error(err))
		return SvcRslt[[]model.WorksetInfo]{}, DB_FAILURE
	}

	infos := make([]model.WorksetInfo, 0, len(details))

	for _, d := range details {
		infos = append(infos, model.WorksetInfo{
			ID:              d.ID,
			Index:           d.Index,
			Name:            d.Name,
			ComicCount:      d.ComicCount,
			Description:     d.Description,
			CreatorID:       d.CreatorID,
			CreatorNickname: d.CreatorNickname,
			CreatedAt:       d.CreatedAt.Unix(),
			UpdatedAt:       d.UpdatedAt.Unix(),
		})
	}

	return accept(200, infos), NO_ERROR
}

func (ws *worksetSvc) CreateWorkset(
	opID string,
	args *model.CreateWorksetArgs,
) (SvcRslt[model.CreateWorksetReply], SvcErr) {
	basicUser, err := ws.userRepo.GetUserByID(nil, opID)
	if err != nil {
		zap.L().Error("Failed to get user info for workset creation", zap.String("userID", opID), zap.Error(err))
		return SvcRslt[model.CreateWorksetReply]{}, DB_FAILURE
	}

	if !basicUser.IsAdmin {
		zap.L().Warn("Non-admin user attempted to create workset", zap.String("userID", opID))
		return SvcRslt[model.CreateWorksetReply]{}, PERMISSION_DENIED
	}

	id, err := genUUID()
	if err != nil {
		zap.L().Error("Failed to generate id for workset", zap.Error(err))
		return SvcRslt[model.CreateWorksetReply]{}, ID_GEN_FAILURE
	}

	nw := &po.NewWorkset{
		ID:          id,
		Name:        args.Name,
		Description: args.Description,
		CreatorID:   opID,
	}

	if err := ws.repo.CreateWorkset(nil, nw); err != nil {
		zap.L().Error("Failed to create workset", zap.Error(err))
		return SvcRslt[model.CreateWorksetReply]{}, DB_FAILURE
	}

	return accept(201, model.CreateWorksetReply{ID: id}), NO_ERROR
}

func (ws *worksetSvc) UpdateWorksetByID(args *model.UpdateWorksetArgs) SvcErr {
	patch := &po.PatchWorkset{
		ID:          args.ID,
		Description: args.Description,
	}

	if err := ws.repo.UpdateWorksetByID(nil, patch); err != nil {
		zap.L().Error("Failed to update workset", zap.Error(err))
		return DB_FAILURE
	}

	return NO_ERROR
}

func (ws *worksetSvc) DeleteWorksetByID(worksetID string) SvcErr {
	if err := ws.repo.DeleteWorksetByID(nil, worksetID); err != nil {
		if err == repo.REC_NOT_FOUND {
			zap.L().Warn("Workset not found for deletion", zap.String("worksetID", worksetID))
			return NOT_FOUND
		}
		zap.L().Error("Failed to delete workset", zap.String("worksetID", worksetID), zap.Error(err))
		return DB_FAILURE
	}

	return NO_ERROR
}
