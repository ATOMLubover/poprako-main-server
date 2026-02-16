package svc

import (
	"time"

	"poprako-main-server/internal/model"
	"poprako-main-server/internal/model/po"
	"poprako-main-server/internal/repo"

	"go.uber.org/zap"
)

// ComicAsgnSvc defines service operations for comic assignments.
type ComicAsgnSvc interface {
	GetAsgnByID(assignmentID string) (SvcRslt[model.ComicAsgnInfo], SvcErr)
	GetAsgnsByComicID(comicID string, offset, limit int) (SvcRslt[[]model.ComicAsgnInfo], SvcErr)
	GetAsgnsByUserID(userID string, offset, limit int) (SvcRslt[[]model.ComicAsgnInfo], SvcErr)

	CreateAsgn(args model.CreateComicAsgnArgs) (SvcRslt[string], SvcErr)

	UpdateAsgnByID(args model.UpdateComicAsgnArgs) SvcErr

	DeleteAsgnByID(assignmentID string) SvcErr
}

type comicAsgnSvc struct {
	repo     repo.ComicAsgnRepo
	userRepo repo.UserRepo
}

// NewComicAsgnSvc creates a new ComicAsgnSvc. r and userRepo must not be nil.
func NewComicAsgnSvc(r repo.ComicAsgnRepo, userRepo repo.UserRepo) ComicAsgnSvc {
	if r == nil {
		panic("ComicAsgnRepo cannot be nil")
	}
	if userRepo == nil {
		panic("UserRepo cannot be nil")
	}

	return &comicAsgnSvc{repo: r, userRepo: userRepo}
}

// GetAsgnByID retrieves a comic assignment by ID.
func (cas *comicAsgnSvc) GetAsgnByID(assignmentID string) (SvcRslt[model.ComicAsgnInfo], SvcErr) {
	asgn, err := cas.repo.GetAsgnByID(nil, assignmentID)
	if err != nil {
		zap.L().Error("Failed to get assignment by ID", zap.String("assignmentID", assignmentID), zap.Error(err))
		return SvcRslt[model.ComicAsgnInfo]{}, DB_FAILURE
	}

	asgnInfo := poAsgnToModelAsgn(asgn)
	return accept(200, asgnInfo), NO_ERROR
}

// GetAsgnsByComicID retrieves comic assignments by comic ID with pagination.
func (cas *comicAsgnSvc) GetAsgnsByComicID(comicID string, offset, limit int) (SvcRslt[[]model.ComicAsgnInfo], SvcErr) {
	asgnList, err := cas.repo.GetAsgnsByComicID(nil, comicID, offset, limit)
	if err != nil {
		zap.L().Error("Failed to get assignments by comic ID", zap.String("comicID", comicID), zap.Error(err))
		return SvcRslt[[]model.ComicAsgnInfo]{}, DB_FAILURE
	}

	asgnInfos := make([]model.ComicAsgnInfo, 0, len(asgnList))
	for _, asgn := range asgnList {
		asgnInfos = append(asgnInfos, poAsgnToModelAsgn(&asgn))
	}

	return accept(200, asgnInfos), NO_ERROR
}

// GetAsgnsByUserID retrieves comic assignments by user ID with pagination.
func (cas *comicAsgnSvc) GetAsgnsByUserID(userID string, offset, limit int) (SvcRslt[[]model.ComicAsgnInfo], SvcErr) {
	asgnList, err := cas.repo.GetAsgnsByUserID(nil, userID, offset, limit)
	if err != nil {
		zap.L().Error("Failed to get assignments by user ID", zap.String("userID", userID), zap.Error(err))
		return SvcRslt[[]model.ComicAsgnInfo]{}, DB_FAILURE
	}

	asgnInfos := make([]model.ComicAsgnInfo, 0, len(asgnList))
	for _, asgn := range asgnList {
		asgnInfos = append(asgnInfos, poAsgnToModelAsgn(&asgn))
	}

	return accept(200, asgnInfos), NO_ERROR
}

// CreateAsgn creates a new comic assignment.
func (cas *comicAsgnSvc) CreateAsgn(args model.CreateComicAsgnArgs) (SvcRslt[string], SvcErr) {
	// Validate user qualifications for requested roles
	user, err := cas.userRepo.GetUserByID(nil, args.AssigneeID)
	if err != nil {
		zap.L().Error("Failed to get user info for assignment", zap.String("userID", args.AssigneeID), zap.Error(err))
		return SvcRslt[string]{}, DB_FAILURE
	}

	// Check if user has required qualifications for requested roles
	if args.IsTranslator != nil && *args.IsTranslator {
		if user.AssignedTranslatorAt == nil {
			zap.L().Warn("User does not have translator qualification", zap.String("userID", args.AssigneeID))
			return SvcRslt[string]{}, PERMISSION_DENIED
		}
	}
	if args.IsProofreader != nil && *args.IsProofreader {
		if user.AssignedProofreaderAt == nil {
			zap.L().Warn("User does not have proofreader qualification", zap.String("userID", args.AssigneeID))
			return SvcRslt[string]{}, PERMISSION_DENIED
		}
	}
	if args.IsTypesetter != nil && *args.IsTypesetter {
		if user.AssignedTypesetterAt == nil {
			zap.L().Warn("User does not have typesetter qualification", zap.String("userID", args.AssigneeID))
			return SvcRslt[string]{}, PERMISSION_DENIED
		}
	}
	if args.IsRedrawer != nil && *args.IsRedrawer {
		if user.AssignedRedrawerAt == nil {
			zap.L().Warn("User does not have redrawer qualification", zap.String("userID", args.AssigneeID))
			return SvcRslt[string]{}, PERMISSION_DENIED
		}
	}
	if args.IsReviewer != nil && *args.IsReviewer {
		if user.AssignedReviewerAt == nil {
			zap.L().Warn("User does not have reviewer qualification", zap.String("userID", args.AssigneeID))
			return SvcRslt[string]{}, PERMISSION_DENIED
		}
	}

	// Generate ID for the assignment
	id, err := genUUID()
	if err != nil {
		zap.L().Error("Failed to generate UUID for assignment", zap.Error(err))
		return SvcRslt[string]{}, ID_GEN_FAILURE
	}

	newAssign := &po.NewComicAsgn{
		ID:      id,
		ComicID: args.ComicID,
		UserID:  args.AssigneeID,
	}

	if err := cas.repo.CreateAsgn(nil, newAssign); err != nil {
		zap.L().Error("Failed to create assignment", zap.Error(err))
		return SvcRslt[string]{}, DB_FAILURE
	}

	// Set role timestamps based on requested roles
	now := time.Now()
	patchAssign := po.PatchComicAsgn{
		ID: id,
	}

	if args.IsTranslator != nil && *args.IsTranslator {
		patchAssign.AssignedTranslatorAt = &now
	}
	if args.IsProofreader != nil && *args.IsProofreader {
		patchAssign.AssignedProofreaderAt = &now
	}
	if args.IsTypesetter != nil && *args.IsTypesetter {
		patchAssign.AssignedTypesetterAt = &now
	}
	if args.IsRedrawer != nil && *args.IsRedrawer {
		patchAssign.AssignedRedrawerAt = &now
	}
	if args.IsReviewer != nil && *args.IsReviewer {
		patchAssign.AssignedReviewerAt = &now
	}

	// Update assignment with role timestamps
	if err := cas.repo.UpdateAsgnByID(nil, &patchAssign); err != nil {
		zap.L().Error("Failed to set assignment roles", zap.Error(err))
		return SvcRslt[string]{}, DB_FAILURE
	}

	return accept(201, id), NO_ERROR
}

// UpdateAsgnByID updates a comic assignment by ID.
func (cas *comicAsgnSvc) UpdateAsgnByID(args model.UpdateComicAsgnArgs) SvcErr {
	patchAssign := modelAsgnArgsToPoPatch(args)

	if err := cas.repo.UpdateAsgnByID(nil, &patchAssign); err != nil {
		zap.L().Error("Failed to update assignment", zap.Error(err))
		return DB_FAILURE
	}

	return NO_ERROR
}

func (cas *comicAsgnSvc) DeleteAsgnByID(assignmentID string) SvcErr {
	if err := cas.repo.DeleteAsgnByID(nil, assignmentID); err != nil {
		if err == repo.REC_NOT_FOUND {
			zap.L().Warn("Assignment not found for deletion", zap.String("assignmentID", assignmentID))
			return NOT_FOUND
		}
		zap.L().Error("Failed to delete assignment", zap.String("assignmentID", assignmentID), zap.Error(err))
		return DB_FAILURE
	}

	return NO_ERROR
}

// poAsgnToModelAsgn converts po.BasicComicAsgn to model.ComicAsgnInfo
func poAsgnToModelAsgn(asgn *po.BasicComicAsgn) model.ComicAsgnInfo {
	info := model.ComicAsgnInfo{
		ID:           asgn.ID,
		ComicID:      asgn.ComicID,
		UserID:       asgn.UserID,
		UserNickname: asgn.UserNickname,
		CreatedAt:    asgn.CreatedAt.Unix(),
		UpdatedAt:    asgn.UpdatedAt.Unix(),
	}

	if asgn.AssignedTranslatorAt != nil {
		ts := asgn.AssignedTranslatorAt.Unix()
		info.AssignedTranslatorAt = &ts
	}
	if asgn.AssignedProofreaderAt != nil {
		ts := asgn.AssignedProofreaderAt.Unix()
		info.AssignedProofreaderAt = &ts
	}
	if asgn.AssignedTypesetterAt != nil {
		ts := asgn.AssignedTypesetterAt.Unix()
		info.AssignedTypesetterAt = &ts
	}
	if asgn.AssignedRedrawerAt != nil {
		ts := asgn.AssignedRedrawerAt.Unix()
		info.AssignedRedrawerAt = &ts
	}
	if asgn.AssignedReviewerAt != nil {
		ts := asgn.AssignedReviewerAt.Unix()
		info.AssignedReviewerAt = &ts
	}

	return info
}

// modelAsgnArgsToPoPatch converts model.UpdateComicAsgnArgs to po.PatchComicAsgn
func modelAsgnArgsToPoPatch(args model.UpdateComicAsgnArgs) po.PatchComicAsgn {
	patch := po.PatchComicAsgn{
		ID: args.ID,
	}

	now := time.Now()
	zeroTime := time.Time{}

	// Convert role flags to timestamp assignments/removals
	// true = assign role with current timestamp
	// false = remove role (set to zero time, which repo converts to NULL)
	if args.IsTranslator != nil {
		if *args.IsTranslator {
			patch.AssignedTranslatorAt = &now
		} else {
			patch.AssignedTranslatorAt = &zeroTime
		}
	}

	if args.IsProofreader != nil {
		if *args.IsProofreader {
			patch.AssignedProofreaderAt = &now
		} else {
			patch.AssignedProofreaderAt = &zeroTime
		}
	}

	if args.IsTypesetter != nil {
		if *args.IsTypesetter {
			patch.AssignedTypesetterAt = &now
		} else {
			patch.AssignedTypesetterAt = &zeroTime
		}
	}

	if args.IsRedrawer != nil {
		if *args.IsRedrawer {
			patch.AssignedRedrawerAt = &now
		} else {
			patch.AssignedRedrawerAt = &zeroTime
		}
	}

	if args.IsReviewer != nil {
		if *args.IsReviewer {
			patch.AssignedReviewerAt = &now
		} else {
			patch.AssignedReviewerAt = &zeroTime
		}
	}

	return patch
}
