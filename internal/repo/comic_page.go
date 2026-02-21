package repo

import (
	"errors"

	"poprako-main-server/internal/model/po"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ComicPageRepo interface {
	Repo

	GetPageByID(ex Exct, pageID string) (*po.BasicComicPage, error)
	GetCoverByComicID(ex Exct, comicID string) (*po.BasicComicPage, error)
	GetPagesByComicID(ex Exct, comicID string) ([]po.BasicComicPage, error)

	CreatePages(newPages []po.NewComicPage) error

	UpdatePageByID(ex Exct, patchPage *po.PatchComicPage) error

	DeletePageByID(ex Exct, pageID string) error
}

type comicPageRepo struct {
	ex Exct
}

func NewComicPageRepo(ex Exct) ComicPageRepo {
	return &comicPageRepo{ex: ex}
}

func (cpr *comicPageRepo) Exct() Exct { return cpr.ex }

func (cpr *comicPageRepo) withTrx(tx Exct) Exct {
	if tx != nil {
		return tx
	}

	return cpr.ex
}

func (cpr *comicPageRepo) CreatePages(newPages []po.NewComicPage) error {
	if err := cpr.Exct().Transaction(func(ex Exct) error {
		cnt := len(newPages)
		if cnt == 0 {
			return nil
		}

		// Insert all pages.
		if err := ex.Create(&newPages).Error; err != nil {
			return err
		}

		// Update comic page_count.
		if err := ex.Model(&po.BasicComic{}).
			Where("id = ?", newPages[0].ComicID).
			UpdateColumn("page_count", gorm.Expr("page_count + ?", cnt)).
			Error; err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (cpr *comicPageRepo) GetPageByID(ex Exct, pageID string) (*po.BasicComicPage, error) {
	ex = cpr.withTrx(ex)

	p := &po.BasicComicPage{}

	if err := ex.
		Where("id = ?", pageID).
		First(p).
		Error; err != nil {
		return nil, err
	}

	return p, nil
}

func (cpr *comicPageRepo) GetPagesByComicID(ex Exct, comicID string) ([]po.BasicComicPage, error) {
	ex = cpr.withTrx(ex)

	var lst []po.BasicComicPage

	if err := ex.
		Where("comic_id = ?", comicID).
		Find(&lst).
		Error; err != nil {
		return nil, err
	}

	return lst, nil
}

func (cpr *comicPageRepo) GetCoverByComicID(ex Exct, comicID string) (*po.BasicComicPage, error) {
	ex = cpr.withTrx(ex)

	p := &po.BasicComicPage{}

	if err := ex.
		Where("comic_id = ?", comicID).
		Order("index ASC").
		First(p).
		Error; err != nil {
		return nil, err
	}

	return p, nil
}

func (cpr *comicPageRepo) UpdatePageByID(ex Exct, patchPage *po.PatchComicPage) error {
	if patchPage.ID == "" {
		return errors.New("page ID is required for update")
	}

	ex = cpr.withTrx(ex)

	updates := map[string]any{}

	if patchPage.ComicID != nil {
		updates["comic_id"] = *patchPage.ComicID
	}
	if patchPage.Index != nil {
		updates["index"] = *patchPage.Index
	}
	if patchPage.OSSKey != nil {
		updates["oss_key"] = *patchPage.OSSKey
	}
	if patchPage.Uploaded != nil {
		updates["uploaded"] = *patchPage.Uploaded
	}

	if len(updates) == 0 {
		return nil
	}

	return ex.Model(&po.PatchComicPage{}).
		Where("id = ?", patchPage.ID).
		Updates(updates).
		Error
}

func (cpr *comicPageRepo) DeletePageByID(ex Exct, pageID string) error {
	ex = cpr.withTrx(ex)

	return ex.Transaction(func(tx Exct) error {
		// Get page first to get comic_id
		page := &po.BasicComicPage{}
		if err := tx.Where("id = ?", pageID).First(page).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return REC_NOT_FOUND
			}
			return err
		}

		comicID := page.ComicID

		// Lock the comic to serialize access to pages
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("id = ?", comicID).
			First(&po.BasicComic{}).Error; err != nil {
			return err
		}

		// Delete the page
		if err := tx.Where("id = ?", pageID).Delete(&po.BasicComicPage{}).Error; err != nil {
			return err
		}

		// Update comic page_count
		if err := tx.Model(&po.BasicComic{}).
			Where("id = ?", comicID).
			UpdateColumn("page_count", gorm.Expr("page_count - ?", 1)).
			Error; err != nil {
			return err
		}

		// Get all pages for the comic to reorder
		var pages []po.BasicComicPage
		if err := tx.Select("id", "index").
			Where("comic_id = ?", comicID).
			Order("index asc").
			Find(&pages).Error; err != nil {
			return err
		}

		// Rewrite indices to ensure continuity
		for i, p := range pages {
			if p.Index != int64(i) {
				if err := tx.Model(&po.BasicComicPage{}).
					Where("id = ?", p.ID).
					UpdateColumn("index", int64(i)).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}
