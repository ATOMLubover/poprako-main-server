package svc

import (
	"fmt"
	"time"

	"poprako-main-server/internal/model"
	"poprako-main-server/internal/model/po"
	"poprako-main-server/internal/repo"

	"go.uber.org/zap"
)

type ComicSvc interface {
	GetComicInfoByID(comicID string) (SvcRslt[model.ComicInfo], SvcErr)
	GetComicBriefsByWorksetID(worksetID string, offset, limit int) (SvcRslt[[]model.ComicBrief], SvcErr)
	RetrieveComics(opt model.RetrieveComicOpt) (SvcRslt[[]model.ComicBrief], SvcErr)

	CreateComic(opID string, args model.CreateComicArgs) (SvcRslt[model.CreateComicReply], SvcErr)

	UpdateComicByID(args model.UpdateComicArgs) SvcErr

	DeleteComicByID(comicID string) SvcErr
}

type comicSvc struct {
	repo          repo.ComicRepo
	userRepo      repo.UserRepo
	comicAsgnRepo repo.ComicAsgnRepo
}

func NewComicSvc(r repo.ComicRepo, ur repo.UserRepo, car repo.ComicAsgnRepo) ComicSvc {
	if r == nil {
		panic("ComicRepo cannot be nil")
	}
	if ur == nil {
		panic("UserRepo cannot be nil")
	}
	if car == nil {
		panic("ComicAsgnRepo cannot be nil")
	}

	return &comicSvc{
		repo:          r,
		userRepo:      ur,
		comicAsgnRepo: car,
	}
}

// GetComicInfoByID retrieves detailed comic info by ID.
func (cs *comicSvc) GetComicInfoByID(comicID string) (SvcRslt[model.ComicInfo], SvcErr) {
	basic, err := cs.repo.GetComicByID(nil, comicID)
	if err != nil {
		zap.L().Error("Failed to get comic by ID", zap.String("comicID", comicID), zap.Error(err))
		return SvcRslt[model.ComicInfo]{}, DB_FAILURE
	}

	info := model.ComicInfo{
		ID:              basic.ID,
		WorksetID:       basic.WorksetID,
		WorksetIndex:    basic.WorksetIndex,
		Index:           basic.Index,
		CreatorID:       basic.CreatorID,
		CreatorNickname: basic.CreatorNickname,
		Author:          basic.Author,
		Title:           basic.Title,
		Description:     basic.Description,
		Comment:         basic.Comment,
		PageCount:       basic.PageCount,
		CreatedAt:       basic.CreatedAt.Unix(),
		UpdatedAt:       basic.UpdatedAt.Unix(),
	}

	// Handle optional timestamp fields
	info.TranslatingStartedAt = timePtrToInt64Ptr(basic.TranslatingStartedAt)
	info.TranslatingCompletedAt = timePtrToInt64Ptr(basic.TranslatingCompletedAt)
	info.ProofreadingStartedAt = timePtrToInt64Ptr(basic.ProofreadingStartedAt)
	info.ProofreadingCompletedAt = timePtrToInt64Ptr(basic.ProofreadingCompletedAt)
	info.TypesettingStartedAt = timePtrToInt64Ptr(basic.TypesettingStartedAt)
	info.TypesettingCompletedAt = timePtrToInt64Ptr(basic.TypesettingCompletedAt)
	info.ReviewingCompletedAt = timePtrToInt64Ptr(basic.ReviewingCompletedAt)
	info.UploadingCompletedAt = timePtrToInt64Ptr(basic.UploadingCompletedAt)

	return accept(200, info), NO_ERROR
}

// GetComicBriefsByWorksetID retrieves brief comic info by workset ID with pagination.
func (cs *comicSvc) GetComicBriefsByWorksetID(worksetID string, offset, limit int) (SvcRslt[[]model.ComicBrief], SvcErr) {
	briefs, err := cs.repo.GetComicsByWorksetID(nil, worksetID, offset, limit)
	if err != nil {
		zap.L().Error("Failed to get comics by workset ID", zap.String("worksetID", worksetID), zap.Error(err))
		return SvcRslt[[]model.ComicBrief]{}, DB_FAILURE
	}

	lst := make([]model.ComicBrief, 0, len(briefs))

	for _, cb := range briefs {
		brief := model.ComicBrief{
			ID:           cb.ID,
			WorksetID:    cb.WorksetID,
			WorksetIndex: cb.WorksetIndex,
			Index:        cb.Index,
			Author:       cb.Author,
			Title:        cb.Title,
			PageCount:    cb.PageCount,
		}

		// Handle optional timestamp fields
		brief.TranslatingStartedAt = timePtrToInt64Ptr(cb.TranslatingStartedAt)
		brief.TranslatingCompletedAt = timePtrToInt64Ptr(cb.TranslatingCompletedAt)
		brief.ProofreadingStartedAt = timePtrToInt64Ptr(cb.ProofreadingStartedAt)
		brief.ProofreadingCompletedAt = timePtrToInt64Ptr(cb.ProofreadingCompletedAt)
		brief.TypesettingStartedAt = timePtrToInt64Ptr(cb.TypesettingStartedAt)
		brief.TypesettingCompletedAt = timePtrToInt64Ptr(cb.TypesettingCompletedAt)
		brief.ReviewingCompletedAt = timePtrToInt64Ptr(cb.ReviewingCompletedAt)
		brief.UploadingCompletedAt = timePtrToInt64Ptr(cb.UploadingCompletedAt)
		lst = append(lst, brief)
	}

	return accept(200, lst), NO_ERROR
}

// RetrieveComics retrieves comics with filtering and pagination.
func (cs *comicSvc) RetrieveComics(opt model.RetrieveComicOpt) (SvcRslt[[]model.ComicBrief], SvcErr) {
	briefs, err := cs.repo.RetrieveComics(nil, opt)
	if err != nil {
		zap.L().Error("Failed to retrieve comics", zap.Error(err))
		return SvcRslt[[]model.ComicBrief]{}, DB_FAILURE
	}

	lst := make([]model.ComicBrief, 0, len(briefs))
	for _, cb := range briefs {
		brief := model.ComicBrief{
			ID:           cb.ID,
			WorksetID:    cb.WorksetID,
			WorksetIndex: cb.WorksetIndex,
			Index:        cb.Index,
			Author:       cb.Author,
			Title:        cb.Title,
			PageCount:    cb.PageCount,
		}

		// Handle optional timestamp fields
		brief.TranslatingStartedAt = timePtrToInt64Ptr(cb.TranslatingStartedAt)
		brief.TranslatingCompletedAt = timePtrToInt64Ptr(cb.TranslatingCompletedAt)
		brief.ProofreadingStartedAt = timePtrToInt64Ptr(cb.ProofreadingStartedAt)
		brief.ProofreadingCompletedAt = timePtrToInt64Ptr(cb.ProofreadingCompletedAt)
		brief.TypesettingStartedAt = timePtrToInt64Ptr(cb.TypesettingStartedAt)
		brief.TypesettingCompletedAt = timePtrToInt64Ptr(cb.TypesettingCompletedAt)
		brief.ReviewingCompletedAt = timePtrToInt64Ptr(cb.ReviewingCompletedAt)
		brief.UploadingCompletedAt = timePtrToInt64Ptr(cb.UploadingCompletedAt)

		lst = append(lst, brief)
	}

	return accept(200, lst), NO_ERROR
}

// CreateComic creates a new comic.
func (cs *comicSvc) CreateComic(opID string, args model.CreateComicArgs) (SvcRslt[model.CreateComicReply], SvcErr) {
	// Check if creator is admin
	creator, err := cs.userRepo.GetUserByID(nil, opID)
	if err != nil {
		zap.L().Error("Failed to get creator info for comic creation", zap.String("userID", opID), zap.Error(err))
		return SvcRslt[model.CreateComicReply]{}, DB_FAILURE
	}

	if !creator.IsAdmin {
		zap.L().Warn("Non-admin user attempted to create comic", zap.String("userID", opID))
		return SvcRslt[model.CreateComicReply]{}, PERMISSION_DENIED
	}

	// Validate pre-assignments
	if len(args.PreAsgns) > 0 {
		if svcErr := cs.validatePreAssignments(args.PreAsgns); svcErr != NO_ERROR {
			return SvcRslt[model.CreateComicReply]{}, svcErr
		}
	}

	// Generate UUID for the new comic
	newID, err := genUUID()
	if err != nil {
		zap.L().Error("Failed to generate UUID for new comic", zap.Error(err))
		return SvcRslt[model.CreateComicReply]{}, ID_GEN_FAILURE
	}

	// Create comic and assignments in a transaction
	if err := cs.repo.Exec().Transaction(func(tx repo.Executor) error {
		// Create the comic
		newComic := &po.NewComic{
			ID:          newID,
			WorksetID:   args.WorksetID,
			CreatorID:   opID,
			Author:      args.Author,
			Title:       args.Title,
			Description: args.Description,
			Comment:     args.Comment,
		}

		if err := cs.repo.CreateComic(newComic); err != nil {
			return fmt.Errorf("failed to create comic: %w", err)
		}

		// Create pre-assignments
		if err := cs.createPreAssignments(tx, newID, args.PreAsgns); err != nil {
			return err
		}

		return nil
	}); err != nil {
		zap.L().Error("Failed to create comic with assignments", zap.String("worksetID", args.WorksetID), zap.Error(err))
		return SvcRslt[model.CreateComicReply]{}, DB_FAILURE
	}

	return accept(201, model.CreateComicReply{ID: newID}), NO_ERROR
}

// UpdateComicByID updates comic info by ID.
func (cs *comicSvc) UpdateComicByID(args model.UpdateComicArgs) SvcErr {
	now := time.Now()

	patch := &po.PatchComic{
		ID:          args.ID,
		Author:      args.Author,
		Title:       args.Title,
		Description: args.Description,
		Comment:     args.Comment,
	}

	// Handle workflow timestamp toggles - only set when true
	if args.TranslatingStarted != nil && *args.TranslatingStarted {
		patch.TranslatingStartedAt = &now
	}
	if args.TranslatingCompleted != nil && *args.TranslatingCompleted {
		patch.TranslatingCompletedAt = &now
	}
	if args.ProofreadingStarted != nil && *args.ProofreadingStarted {
		patch.ProofreadingStartedAt = &now
	}
	if args.ProofreadingCompleted != nil && *args.ProofreadingCompleted {
		patch.ProofreadingCompletedAt = &now
	}
	if args.TypesettingStarted != nil && *args.TypesettingStarted {
		patch.TypesettingStartedAt = &now
	}
	if args.TypesettingCompleted != nil && *args.TypesettingCompleted {
		patch.TypesettingCompletedAt = &now
	}
	if args.ReviewingCompleted != nil && *args.ReviewingCompleted {
		patch.ReviewingCompletedAt = &now
	}
	if args.UploadingCompleted != nil && *args.UploadingCompleted {
		patch.UploadingCompletedAt = &now
	}

	if err := cs.repo.UpdateComicByID(nil, patch); err != nil {
		zap.L().Error("Failed to update comic", zap.String("comicID", args.ID), zap.Error(err))
		return DB_FAILURE
	}

	return NO_ERROR
}

// validatePreAssignments validates that all pre-assigned users have the required qualifications.
func (cs *comicSvc) validatePreAssignments(preAsgns []model.PreAsgnArgs) SvcErr {
	for _, preAsgn := range preAsgns {
		user, err := cs.userRepo.GetUserByID(nil, preAsgn.AssigneeID)
		if err != nil {
			zap.L().Error("Failed to get user info for pre-assignment validation",
				zap.String("userID", preAsgn.AssigneeID), zap.Error(err))
			return DB_FAILURE
		}

		// Check if user has required qualifications
		if preAsgn.IsTranslator != nil && *preAsgn.IsTranslator {
			if user.AssignedTranslatorAt == nil {
				zap.L().Warn("User does not have translator qualification",
					zap.String("userID", preAsgn.AssigneeID))
				return PERMISSION_DENIED
			}
		}

		if preAsgn.IsProofreader != nil && *preAsgn.IsProofreader {
			if user.AssignedProofreaderAt == nil {
				zap.L().Warn("User does not have proofreader qualification",
					zap.String("userID", preAsgn.AssigneeID))
				return PERMISSION_DENIED
			}
		}

		if preAsgn.IsTypesetter != nil && *preAsgn.IsTypesetter {
			if user.AssignedTypesetterAt == nil {
				zap.L().Warn("User does not have typesetter qualification",
					zap.String("userID", preAsgn.AssigneeID))
				return PERMISSION_DENIED
			}
		}

		if preAsgn.IsRedrawer != nil && *preAsgn.IsRedrawer {
			if user.AssignedRedrawerAt == nil {
				zap.L().Warn("User does not have redrawer qualification",
					zap.String("userID", preAsgn.AssigneeID))
				return PERMISSION_DENIED
			}
		}

		if preAsgn.IsReviewer != nil && *preAsgn.IsReviewer {
			if user.AssignedReviewerAt == nil {
				zap.L().Warn("User does not have reviewer qualification",
					zap.String("userID", preAsgn.AssigneeID))
				return PERMISSION_DENIED
			}
		}
	}

	return NO_ERROR
}

// createPreAssignments creates comic assignments for pre-assigned users.
// For each pre-assignment, it creates a new assignment record and sets the role timestamps.
func (cs *comicSvc) createPreAssignments(tx repo.Executor, comicID string, preAsgns []model.PreAsgnArgs) error {
	now := time.Now()

	for _, preAsgn := range preAsgns {
		// Generate assignment ID
		asgnID, err := genUUID()
		if err != nil {
			return fmt.Errorf("failed to generate assignment ID: %w", err)
		}

		// Create the assignment record
		newAsgn := &po.NewComicAsgn{
			ID:      asgnID,
			ComicID: comicID,
			UserID:  preAsgn.AssigneeID,
		}

		if err := cs.comicAsgnRepo.CreateAsgn(tx, newAsgn); err != nil {
			return fmt.Errorf("failed to create assignment: %w", err)
		}

		// Set role timestamps based on pre-assignment roles
		patchAsgn := &po.PatchComicAsgn{
			ID: asgnID,
		}

		if preAsgn.IsTranslator != nil && *preAsgn.IsTranslator {
			patchAsgn.AssignedTranslatorAt = &now
		}
		if preAsgn.IsProofreader != nil && *preAsgn.IsProofreader {
			patchAsgn.AssignedProofreaderAt = &now
		}
		if preAsgn.IsTypesetter != nil && *preAsgn.IsTypesetter {
			patchAsgn.AssignedTypesetterAt = &now
		}
		if preAsgn.IsRedrawer != nil && *preAsgn.IsRedrawer {
			patchAsgn.AssignedRedrawerAt = &now
		}
		if preAsgn.IsReviewer != nil && *preAsgn.IsReviewer {
			patchAsgn.AssignedReviewerAt = &now
		}

		// Update the assignment with role timestamps
		if err := cs.comicAsgnRepo.UpdateAsgnByID(tx, patchAsgn); err != nil {
			return fmt.Errorf("failed to update assignment roles: %w", err)
		}
	}

	return nil
}

func (cs *comicSvc) DeleteComicByID(comicID string) SvcErr {
	if err := cs.repo.DeleteComicByID(nil, comicID); err != nil {
		if err == repo.REC_NOT_FOUND {
			zap.L().Warn("Comic not found for deletion", zap.String("comicID", comicID))
			return NOT_FOUND
		}
		zap.L().Error("Failed to delete comic", zap.String("comicID", comicID), zap.Error(err))
		return DB_FAILURE
	}

	return NO_ERROR
}
