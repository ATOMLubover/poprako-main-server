package seeder

import (
	"fmt"
	"math/rand"
	"time"

	"poprako-main-server/internal/model/po"
	"poprako-main-server/internal/repo"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ComicProgress struct {
	ComicID         string
	Pages           []po.NewComicPage
	HasTranslator   bool
	HasProofreader  bool
	HasTypesetter   bool
	HasReviewer     bool
	TranslatorID    string
	ProofreaderID   string
	TranslationRate float64
	ProofreadRate   float64
	TypesettingDone bool
	ReviewingDone   bool
}

func Seed(ex repo.Executor) {
	var count int64
	if err := ex.Model(&po.BriefComic{}).Count(&count).Error; err != nil {
		zap.L().Error("Failed to count comics", zap.Error(err))
		return
	}

	if count > 0 {
		zap.L().Info("Database already has comic data. Skipping seed.", zap.Int64("count", count))
		return
	}

	zap.L().Info("Starting database seeding...")

	userRepo := repo.NewUserRepo(ex)
	worksetRepo := repo.NewWorksetRepo(ex)
	comicRepo := repo.NewComicRepo(ex)
	pageRepo := repo.NewComicPageRepo(ex)
	unitRepo := repo.NewComicUnitRepo(ex)
	asgnRepo := repo.NewComicAsgnRepo(ex)

	rand.Seed(time.Now().UnixNano())

	// 1. Get admin user
	adminUser := &po.BasicUser{}
	if err := ex.Model(&po.BasicUser{}).First(adminUser).Error; err != nil {
		zap.L().Error("Failed to find admin user for seeding", zap.Error(err))
		return
	}

	// 2. Create mock users with different roles
	translators := createUsers(ex, userRepo, 3, "Translator", true, false, false, false, false, false)
	proofreaders := createUsers(ex, userRepo, 3, "Proofreader", false, true, false, false, false, false)
	_ = createUsers(ex, userRepo, 2, "Typesetter", false, false, true, false, false, false)
	_ = createUsers(ex, userRepo, 1, "Reviewer", false, false, false, false, true, false)
	multiRole := createUsers(ex, userRepo, 2, "Multi", true, true, false, false, false, false)

	allTranslators := append(translators, multiRole...)
	allProofreaders := append(proofreaders, multiRole...)

	// 3. Create worksets
	worksets := []string{}
	worksetNames := []string{"连载系列A", "连载系列B", "短篇合集", "特别企划"}
	for i, name := range worksetNames {
		wsID := uuid.NewString()
		ws := &po.NewWorkset{
			ID:          wsID,
			Name:        name,
			CreatorID:   adminUser.ID,
			Description: ptr(fmt.Sprintf("Mock workset %d", i+1)),
		}
		if err := worksetRepo.CreateWorkset(nil, ws); err != nil {
			zap.L().Error("Failed to create workset", zap.Error(err))
			continue
		}
		worksets = append(worksets, wsID)
	}

	// 4. Create comics and track progress
	totalComics := 100
	comicProgresses := make([]ComicProgress, 0, totalComics)

	for i := 0; i < totalComics; i++ {
		comicID := uuid.NewString()
		wsID := worksets[i%len(worksets)]

		comic := &po.NewComic{
			ID:          comicID,
			WorksetID:   wsID,
			CreatorID:   adminUser.ID,
			Author:      fmt.Sprintf("作者%d", (i%20)+1),
			Title:       fmt.Sprintf("漫画作品 #%d", i+1),
			Description: ptr(fmt.Sprintf("测试用漫画作品 %d", i+1)),
		}

		// Note: workset_index will be populated by CreateComic from the repo layer
		if err := comicRepo.CreateComic(comic); err != nil {
			zap.L().Error("Failed to create comic during seeding", zap.Error(err))
			continue
		}

		// Generate pages
		pageCount := rand.Intn(101) + 20
		newPages := make([]po.NewComicPage, 0, pageCount)

		for p := 0; p < pageCount; p++ {
			pageID := uuid.NewString()
			uploaded := rand.Float64() < 0.85
			newPages = append(newPages, po.NewComicPage{
				ID:       pageID,
				ComicID:  comicID,
				Index:    int64(p + 1),
				OSSKey:   fmt.Sprintf("mock/%s/page_%d.jpg", comicID, p+1),
				Uploaded: &uploaded,
			})
		}

		if err := pageRepo.CreatePages(newPages); err != nil {
			zap.L().Error("Failed to create pages during seeding", zap.Error(err))
			continue
		}

		// Determine assignments
		hasTranslator := rand.Float64() < 0.75
		hasProofreader := hasTranslator && rand.Float64() < 0.6
		hasTypesetter := hasProofreader && rand.Float64() < 0.4
		hasReviewer := hasTypesetter && rand.Float64() < 0.3

		progress := ComicProgress{
			ComicID:        comicID,
			Pages:          newPages,
			HasTranslator:  hasTranslator,
			HasProofreader: hasProofreader,
			HasTypesetter:  hasTypesetter,
			HasReviewer:    hasReviewer,
		}

		if hasTranslator {
			progress.TranslatorID = allTranslators[rand.Intn(len(allTranslators))]
			progress.TranslationRate = 0.5 + rand.Float64()*0.5
		}

		if hasProofreader {
			progress.ProofreaderID = allProofreaders[rand.Intn(len(allProofreaders))]
			progress.ProofreadRate = rand.Float64() * 0.9
		}

		if hasTypesetter {
			progress.TypesettingDone = rand.Float64() < 0.5
		}

		if hasReviewer {
			progress.ReviewingDone = rand.Float64() < 0.6
		}

		comicProgresses = append(comicProgresses, progress)
	}

	// 5. Create assignments
	for _, prog := range comicProgresses {
		asgnID := uuid.NewString()
		asgn := &po.NewComicAsgn{
			ID:      asgnID,
			ComicID: prog.ComicID,
			UserID:  adminUser.ID,
		}

		if err := asgnRepo.CreateAsgn(nil, asgn); err != nil {
			zap.L().Error("Failed to create assignment", zap.Error(err))
			continue
		}

		now := time.Now()
		patch := &po.PatchComicAsgn{ID: asgnID}

		if prog.HasTranslator {
			patch.AssignedTranslatorAt = &now
			asgn.UserID = prog.TranslatorID
		}
		if prog.HasProofreader {
			patch.AssignedProofreaderAt = &now
		}
		if prog.HasTypesetter {
			patch.AssignedTypesetterAt = &now
		}
		if prog.HasReviewer {
			patch.AssignedReviewerAt = &now
		}

		if err := asgnRepo.UpdateAsgnByID(nil, patch); err != nil {
			zap.L().Error("Failed to update assignment roles", zap.Error(err))
		}
	}

	// 6. Create units and set progress
	for _, prog := range comicProgresses {
		if !prog.HasTranslator {
			continue
		}

		allUnits := make([]po.NewComicUnit, 0)
		translatedUnits := make([]po.NewComicUnit, 0)

		for _, page := range prog.Pages {
			unitCount := rand.Intn(9) + 3
			for u := 0; u < unitCount; u++ {
				unit := po.NewComicUnit{
					ID:          uuid.NewString(),
					PageID:      page.ID,
					Index:       int64(u + 1),
					XCoordinate: rand.Float64() * 1920,
					YCoordinate: rand.Float64() * 1080,
					IsInBox:     rand.Float64() < 0.65,
					CreatorID:   &prog.TranslatorID,
				}

				if rand.Float64() < prog.TranslationRate {
					txt := fmt.Sprintf("翻译文本 %d", u+1)
					unit.TranslatedText = &txt
					unit.TranslatorID = &prog.TranslatorID
					if rand.Float64() < 0.1 {
						comment := "需要确认"
						unit.TranslatorComment = &comment
					}
					translatedUnits = append(translatedUnits, unit)
				}

				allUnits = append(allUnits, unit)
			}
		}

		if len(allUnits) > 0 {
			if err := unitRepo.CreateUnits(nil, allUnits); err != nil {
				zap.L().Error("Failed to create units", zap.Error(err))
				continue
			}
		}

		// Set proofread status
		if prog.HasProofreader && len(translatedUnits) > 0 {
			proofreadCount := int(float64(len(translatedUnits)) * prog.ProofreadRate)
			for i := 0; i < proofreadCount; i++ {
				unit := &translatedUnits[i]
				proved := true
				provedTxt := *unit.TranslatedText
				if rand.Float64() < 0.25 {
					provedTxt = provedTxt + " [已修正]"
				}

				patch := &po.PatchComicUnit{
					ID:            unit.ID,
					Proved:        &proved,
					ProvedText:    &provedTxt,
					ProofreaderID: &prog.ProofreaderID,
				}

				if rand.Float64() < 0.08 {
					comment := "确认无误"
					patch.ProofreaderComment = &comment
				}

				if err := unitRepo.UpdateUnitsByIDs(nil, []po.PatchComicUnit{*patch}); err != nil {
					zap.L().Error("Failed to update unit proofread status", zap.Error(err))
				}
			}
		}

		// Set comic workflow timestamps
		updateComicProgress(ex, comicRepo, prog, len(allUnits), len(translatedUnits))
	}

	zap.L().Info("Database seeding completed", zap.Int("comics", totalComics))
}

func createUsers(ex repo.Executor, userRepo repo.UserRepo, count int, rolePrefix string,
	translator, proofreader, typesetter, redrawer, reviewer, uploader bool,
) []string {
	now := time.Now()
	users := make([]string, 0, count)

	for i := 0; i < count; i++ {
		userID := uuid.NewString()
		user := &po.NewUser{
			ID:           userID,
			QQ:           fmt.Sprintf("%d%04d", 20000+i, rand.Intn(10000)),
			Nickname:     fmt.Sprintf("%s用户%d", rolePrefix, i+1),
			PasswordHash: "mock_hash",
		}

		if translator {
			user.AssignedTranslatorAt = &now
		}
		if proofreader {
			user.AssignedProofreaderAt = &now
		}
		if typesetter {
			user.AssignedTypesetterAt = &now
		}
		if redrawer {
			user.AssignedRedrawerAt = &now
		}
		if reviewer {
			user.AssignedReviewerAt = &now
		}
		if uploader {
			user.AssignedUploaderAt = &now
		}

		if err := userRepo.CreateUser(nil, user); err != nil {
			zap.L().Error("Failed to create mock user", zap.Error(err))
			continue
		}

		users = append(users, userID)
	}

	return users
}

func updateComicProgress(ex repo.Executor, comicRepo repo.ComicRepo, prog ComicProgress,
	totalUnits, translatedUnits int,
) {
	if !prog.HasTranslator {
		return
	}

	now := time.Now()
	baseTime := now.Add(-time.Hour * 24 * 30)

	patch := &po.PatchComic{ID: prog.ComicID}

	translatingStart := baseTime
	patch.TranslatingStartedAt = &translatingStart

	translationComplete := float64(translatedUnits) / float64(totalUnits)
	if translationComplete > 0.9 {
		translatingEnd := baseTime.Add(time.Hour * 24 * 7)
		patch.TranslatingCompletedAt = &translatingEnd

		if prog.HasProofreader {
			proofStart := translatingEnd.Add(time.Hour * 24)
			patch.ProofreadingStartedAt = &proofStart

			if prog.ProofreadRate > 0.9 {
				proofEnd := proofStart.Add(time.Hour * 24 * 5)
				patch.ProofreadingCompletedAt = &proofEnd

				if prog.HasTypesetter {
					typeStart := proofEnd.Add(time.Hour * 24)
					patch.TypesettingStartedAt = &typeStart

					if prog.TypesettingDone {
						typeEnd := typeStart.Add(time.Hour * 24 * 3)
						patch.TypesettingCompletedAt = &typeEnd

						if prog.HasReviewer && prog.ReviewingDone {
							reviewEnd := typeEnd.Add(time.Hour * 24 * 2)
							patch.ReviewingCompletedAt = &reviewEnd

							if rand.Float64() < 0.5 {
								uploadEnd := reviewEnd.Add(time.Hour * 24)
								patch.UploadingCompletedAt = &uploadEnd
							}
						}
					}
				}
			}
		}
	}

	if err := comicRepo.UpdateComicByID(nil, patch); err != nil {
		zap.L().Error("Failed to update comic progress timestamps", zap.Error(err))
	}
}

func ptr(s string) *string {
	return &s
}

func ptrBool(b bool) *bool {
	return &b
}
