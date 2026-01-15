package svc

import (
	"poprako-main-server/internal/model"
	"poprako-main-server/internal/model/po"
	"poprako-main-server/internal/repo"

	"go.uber.org/zap"
)

// ComicUnitSvc defines service operations for comic units.
type ComicUnitSvc interface {
	GetUnitsByPageID(pageID string) (SvcRslt[[]model.ComicUnitInfo], SvcErr)

	CreateUnits(opID string, newUnits []model.NewComicUnitArgs) SvcErr

	UpdateUnitsByIDs(opID string, patchUnits []model.PatchComicUnitArgs) SvcErr

	DeleteUnitByIDs(unitIDs []string) SvcErr
}

type comicUnitSvc struct {
	repo repo.ComicUnitRepo
}

// NewComicUnitSvc creates a new ComicUnitSvc. r must not be nil.
func NewComicUnitSvc(r repo.ComicUnitRepo) ComicUnitSvc {
	if r == nil {
		panic("ComicUnitRepo cannot be nil")
	}

	return &comicUnitSvc{repo: r}
}

// GetUnitsByPageID retrieves comic units by page ID.
func (cus *comicUnitSvc) GetUnitsByPageID(pageID string) (SvcRslt[[]model.ComicUnitInfo], SvcErr) {
	units, err := cus.repo.GetUnitsByPageID(nil, pageID)
	if err != nil {
		zap.L().Error("Failed to get units by page ID", zap.String("pageID", pageID), zap.Error(err))
		return SvcRslt[[]model.ComicUnitInfo]{}, DB_FAILURE
	}

	// Convert po.BasicComicUnit to model.ComicUnitInfo
	var infos []model.ComicUnitInfo
	for _, u := range units {
		infos = append(infos, model.ComicUnitInfo{
			ID:                 u.ID,
			PageID:             u.PageID,
			Index:              u.Index,
			XCoordinate:        u.XCoordinate,
			YCoordinate:        u.YCoordinate,
			IsInBox:            u.IsInBox,
			TranslatedText:     u.TranslatedText,
			TranslatorID:       u.TranslatorID,
			TranslatorComment:  u.TranslatorComment,
			ProvedText:         u.ProvedText,
			Proved:             u.Proved,
			ProofreaderID:      u.ProofreaderID,
			ProofreaderComment: u.ProofreaderComment,
			CreatorID:          u.CreatorID,
			CreatedAt:          u.CreatedAt.Unix(),
			UpdatedAt:          u.UpdatedAt.Unix(),
		})
	}

	return accept(200, infos), NO_ERROR
}

// CreateUnits creates a batch of comic units.
func (cus *comicUnitSvc) CreateUnits(opID string, newUnits []model.NewComicUnitArgs) SvcErr {
	if len(newUnits) == 0 {
		return NO_ERROR
	}

	// Convert model.NewComicUnitArgs to po.NewComicUnit
	var poUnits []po.NewComicUnit
	for _, u := range newUnits {
		// Generate ID for each unit
		id, err := genUUID()
		if err != nil {
			zap.L().Error("Failed to generate UUID for comic unit", zap.Error(err))
			return ID_GEN_FAILURE
		}

		poUnits = append(poUnits, po.NewComicUnit{
			ID:                 id,
			PageID:             u.PageID,
			Index:              u.Index,
			XCoordinate:        u.XCoordinate,
			YCoordinate:        u.YCoordinate,
			IsInBox:            u.IsInBox,
			TranslatedText:     u.TranslatedText,
			TranslatorID:       nil, // Will be set by update if needed
			TranslatorComment:  u.TranslatorComment,
			ProvedText:         u.ProvedText,
			Proved:             u.Proved,
			ProofreaderID:      nil, // Will be set by update if needed
			ProofreaderComment: u.ProofreaderComment,
			CreatorID:          &opID,
		})
	}

	if err := cus.repo.CreateUnits(nil, poUnits); err != nil {
		zap.L().Error("Failed to create units", zap.Error(err))
		return DB_FAILURE
	}

	return NO_ERROR
}

// UpdateUnitsByIDs updates a batch of comic units by their IDs.
func (cus *comicUnitSvc) UpdateUnitsByIDs(opID string, patchUnits []model.PatchComicUnitArgs) SvcErr {
	if len(patchUnits) == 0 {
		return NO_ERROR
	}

	// Convert model.PatchComicUnitArgs to po.PatchComicUnit
	var poPatches []po.PatchComicUnit
	for _, pu := range patchUnits {
		if pu.ID == "" {
			zap.L().Error("PatchComicUnitArgs missing ID", zap.Any("patchUnit", pu))
			return INVALID_UNIT_DATA
		}

		poPatch := po.PatchComicUnit{
			ID:                 pu.ID,
			Index:              pu.Index,
			XCoordinate:        pu.XCoordinate,
			YCoordinate:        pu.YCoordinate,
			IsInBox:            pu.IsInBox,
			TranslatedText:     pu.TranslatedText,
			TranslatorComment:  pu.TranslatorComment,
			ProvedText:         pu.ProvedText,
			Proved:             pu.Proved,
			ProofreaderComment: pu.ProofreaderComment,
		}

		// Set translator/proofreader ID based on what's being modified
		if pu.TranslatedText != nil {
			poPatch.TranslatorID = &opID
		}
		if pu.ProvedText != nil || pu.Proved != nil {
			poPatch.ProofreaderID = &opID
		}

		poPatches = append(poPatches, poPatch)
	}

	if err := cus.repo.UpdateUnitsByIDs(nil, poPatches); err != nil {
		zap.L().Error("Failed to update units", zap.Error(err))
		return DB_FAILURE
	}

	return NO_ERROR
}

// DeleteUnitByIDs deletes a batch of comic units by their IDs.
func (cus *comicUnitSvc) DeleteUnitByIDs(unitIDs []string) SvcErr {
	if len(unitIDs) == 0 {
		return NO_ERROR
	}

	if err := cus.repo.DeleteUnitByIDs(nil, unitIDs); err != nil {
		zap.L().Error("Failed to delete units", zap.Error(err))
		return DB_FAILURE
	}

	return NO_ERROR
}
