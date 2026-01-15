package repo

import (
	"sync"
	"time"

	"poprako-main-server/internal/model/po"
)

// NewMockWorksetRepo returns an in-memory mock implementing WorksetRepo for tests.
func NewMockWorksetRepo() WorksetRepo {
	return &mockWorksetRepo{store: make(map[string]*po.DetailedWorkset)}
}

type mockWorksetRepo struct {
	mu    sync.Mutex
	store map[string]*po.DetailedWorkset
	ex    Executor
}

func (m *mockWorksetRepo) Exec() Executor { return m.ex }
func (m *mockWorksetRepo) withTrx(tx Executor) Executor {
	if tx != nil {
		return tx
	}
	return m.ex
}

func (m *mockWorksetRepo) GetWorksetByID(ex Executor, worksetID string) (*po.DetailedWorkset, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if w, ok := m.store[worksetID]; ok {
		return w, nil
	}
	return nil, nil
}

func (m *mockWorksetRepo) RetrieveWorksets(ex Executor, limit, offset int) ([]po.DetailedWorkset, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var all []po.DetailedWorkset
	for _, w := range m.store {
		all = append(all, *w)
	}

	if offset < 0 {
		offset = 0
	}
	if offset >= len(all) {
		return []po.DetailedWorkset{}, nil
	}
	start := offset
	end := len(all)
	if limit > 0 && start+limit < end {
		end = start + limit
	}
	return all[start:end], nil
}

func (m *mockWorksetRepo) CreateWorkset(ex Executor, newWorkset *po.NewWorkset) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now()
	var desc *string
	if newWorkset.Description != nil {
		v := *newWorkset.Description
		desc = &v
	}
	m.store[newWorkset.ID] = &po.DetailedWorkset{
		ID:              newWorkset.ID,
		Index:           newWorkset.Index,
		ComicCount:      newWorkset.ComicCount,
		Name:            "",
		Description:     desc,
		CreatorID:       newWorkset.CreatorID,
		CreatorNickname: "",
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	return nil
}

func (m *mockWorksetRepo) UpdateWorksetByID(ex Executor, patchWorkset *po.PatchWorkset) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	w, ok := m.store[patchWorkset.ID]
	if !ok {
		return nil
	}
	if patchWorkset.Index != nil {
		w.Index = *patchWorkset.Index
	}
	if patchWorkset.ComicCount != nil {
		w.ComicCount = *patchWorkset.ComicCount
	}
	if patchWorkset.Description != nil {
		w.Description = patchWorkset.Description
	}
	if patchWorkset.Name != nil {
		w.Name = *patchWorkset.Name
	}
	if patchWorkset.CreatorID != nil {
		w.CreatorID = *patchWorkset.CreatorID
	}
	w.UpdatedAt = time.Now()
	m.store[patchWorkset.ID] = w
	return nil
}

func (m *mockWorksetRepo) DeleteWorksetByID(ex Executor, worksetID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.store[worksetID]; !exists {
		return REC_NOT_FOUND
	}

	delete(m.store, worksetID)
	return nil
}
