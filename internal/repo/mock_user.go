package repo

import (
	"sync"

	"poprako-main-server/internal/model"
	"poprako-main-server/internal/model/po"
)

// NewMockUserRepo returns an in-memory mock implementing UserRepo for tests.
func NewMockUserRepo() UserRepo {
	return &mockUserRepo{basics: make(map[string]*po.BasicUser), byQQ: make(map[string]*po.BasicUser), secrets: make(map[string]*po.SecretUser)}
}

type mockUserRepo struct {
	mu      sync.Mutex
	basics  map[string]*po.BasicUser
	byQQ    map[string]*po.BasicUser
	secrets map[string]*po.SecretUser
	ex      Executor
}

func (m *mockUserRepo) Exec() Executor { return m.ex }
func (m *mockUserRepo) withTrx(tx Executor) Executor {
	if tx != nil {
		return tx
	}

	return m.ex
}

func (m *mockUserRepo) GetUserByID(ex Executor, userID string) (*po.BasicUser, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if b, ok := m.basics[userID]; ok {
		return b, nil
	}
	return nil, nil
}

func (m *mockUserRepo) GetUserByQQ(ex Executor, qq string) (*po.BasicUser, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if b, ok := m.byQQ[qq]; ok {
		return b, nil
	}
	return nil, nil
}

func (m *mockUserRepo) GetSecretUserByQQ(ex Executor, qq string) (*po.SecretUser, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if s, ok := m.secrets[qq]; ok {
		return s, nil
	}
	return nil, nil
}

func (m *mockUserRepo) CreateUser(ex Executor, newUser *po.NewUser) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	b := &po.BasicUser{ID: newUser.ID, QQ: newUser.QQ, Nickname: newUser.Nickname}
	m.basics[newUser.ID] = b
	m.byQQ[newUser.QQ] = b
	m.secrets[newUser.QQ] = &po.SecretUser{ID: newUser.ID, PwdHash: newUser.PasswordHash}
	return nil
}

func (m *mockUserRepo) UpdateUserByID(ex Executor, updateUser *po.PatchUser) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if b, ok := m.basics[updateUser.ID]; ok {
		if updateUser.Nickname != nil {
			b.Nickname = *updateUser.Nickname
		}
		if updateUser.QQ != nil {
			// update byQQ map
			delete(m.byQQ, b.QQ)
			b.QQ = *updateUser.QQ
			m.byQQ[b.QQ] = b
		}
		if updateUser.IsAdmin != nil {
			b.IsAdmin = *updateUser.IsAdmin
		}
		return nil
	}
	return nil
}

func (m *mockUserRepo) RetrieveUsers(ex Executor, opt model.RetrieveUserOpt) ([]po.BasicUser, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var all []po.BasicUser

	for _, b := range m.basics {
		all = append(all, *b)
	}

	offset := opt.Offset
	limit := opt.Limit

	if offset < 0 {
		offset = 0
	}

	if offset >= len(all) {
		return []po.BasicUser{}, nil
	}

	start := offset
	end := len(all)

	if limit > 0 {
		if start+limit < end {
			end = start + limit
		}
	}

	return all[start:end], nil
}
