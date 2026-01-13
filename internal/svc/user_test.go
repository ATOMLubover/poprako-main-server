package svc_test

import (
	"errors"
	"os"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"saas-template-go/internal/jwtcodec"
	"saas-template-go/internal/model"
	"saas-template-go/internal/model/po"
	"saas-template-go/internal/repo"
	svc "saas-template-go/internal/svc"
)

// mockUserRepo implements repo.UserRepo for testing purposes.
type mockUserRepo struct {
	pwdHash   string
	pwdErr    error
	createErr error
	createID  string
}

func (m *mockUserRepo) GetUserBasicByID(ex repo.Executor, userID string) (*po.UserBasic, error) {
	return &po.UserBasic{ID: userID}, nil
}

func (m *mockUserRepo) GetUserBasicByEmail(ex repo.Executor, email string) (*po.UserBasic, error) {
	return &po.UserBasic{Email: email, ID: "id-from-email"}, nil
}

func (m *mockUserRepo) GetUserBasicsByIDs(ex repo.Executor, userIDs []string) ([]po.UserBasic, error) {
	var out []po.UserBasic
	for _, id := range userIDs {
		out = append(out, po.UserBasic{ID: id})
	}
	return out, nil
}

func (m *mockUserRepo) GetUserBasicsByEmails(ex repo.Executor, emails []string) ([]po.UserBasic, error) {
	var out []po.UserBasic
	for _, e := range emails {
		out = append(out, po.UserBasic{Email: e})
	}
	return out, nil
}

func (m *mockUserRepo) GetPwdHashByEmail(ex repo.Executor, email string) (string, error) {
	return m.pwdHash, m.pwdErr
}

func (m *mockUserRepo) CreateNewUser(ex repo.Executor, newUser *po.NewUser) error {
	if m.createErr != nil {
		return m.createErr
	}

	if m.createID != "" {
		newUser.ID = m.createID
	} else {
		newUser.ID = "mock-generated-id"
	}

	return nil
}

func (m *mockUserRepo) UpdateUserByID(ex repo.Executor, updateUser *po.UpdateUser) error {
	if updateUser.ID == "" {
		return errors.New("id required")
	}
	return nil
}

func TestLoginUser(t *testing.T) {
	os.Setenv("JWT_SECRET_KEY", "test-secret")
	defer os.Unsetenv("JWT_SECRET_KEY")

	codec := jwtcodec.NewJWTCodec(3600)

	t.Run("create new user on login", func(t *testing.T) {
		mr := &mockUserRepo{
			pwdHash:  "",
			createID: "new-user-id",
		}

		us := svc.NewUserSvc(mr, codec)

		rslt, err := us.LoginUser(model.LoginArgs{Email: "new@example.com", Password: "pass123"})
		if err != svc.NO_ERROR {
			t.Fatalf("expected NO_ERROR, got %v", err)
		}

		if rslt.Code != 204 {
			t.Fatalf("expected code 204 for new user creation, got %d", rslt.Code)
		}

		if rslt.Data == nil || rslt.Data.Token == "" {
			t.Fatalf("expected token in result, got %+v", rslt)
		}
	})

	t.Run("existing user login success", func(t *testing.T) {
		pass := "secret-pass"
		hash, _ := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)

		mr := &mockUserRepo{pwdHash: string(hash)}

		us := svc.NewUserSvc(mr, codec)

		rslt, err := us.LoginUser(model.LoginArgs{Email: "exist@example.com", Password: pass})
		if err != svc.NO_ERROR {
			t.Fatalf("expected NO_ERROR, got %v", err)
		}
		if rslt.Code != 200 {
			t.Fatalf("expected code 200 for existing user login, got %d", rslt.Code)
		}
		if rslt.Data == nil || rslt.Data.Token == "" {
			t.Fatalf("expected token in result, got %+v", rslt)
		}
	})

	t.Run("password mismatch", func(t *testing.T) {
		// hash for a different password
		hash, _ := bcrypt.GenerateFromPassword([]byte("other-pass"), bcrypt.DefaultCost)
		mr := &mockUserRepo{pwdHash: string(hash)}

		us := svc.NewUserSvc(mr, codec)

		_, err := us.LoginUser(model.LoginArgs{Email: "exist@example.com", Password: "attempt"})
		if err != svc.PWD_MISMATCH {
			t.Fatalf("expected PWD_MISMATCH, got %v", err)
		}
	})

	t.Run("db failure on get pwd", func(t *testing.T) {
		mr := &mockUserRepo{pwdErr: errors.New("db fail")}

		us := svc.NewUserSvc(mr, codec)

		_, err := us.LoginUser(model.LoginArgs{Email: "err@example.com", Password: "p"})
		if err != svc.DB_FAILURE {
			t.Fatalf("expected DB_FAILURE, got %v", err)
		}
	})
}
