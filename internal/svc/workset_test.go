package svc

// import (
// 	"testing"

// 	"poprako-main-server/internal/model"
// 	"poprako-main-server/internal/repo"
// )

// // mockUserSvc implements UserSvc for tests; treats any user as admin.
// type mockUserSvc struct{}

// func (m *mockUserSvc) GetUserInfoByID(userID string) (SvcRslt[model.UserInfo], SvcErr) {
// 	return accept(200, model.UserInfo{UserID: userID, IsAdmin: true}), NO_ERROR
// }

// func (m *mockUserSvc) GetUserInfoByQQ(qq string) (SvcRslt[model.UserInfo], SvcErr) {
// 	return SvcRslt[model.UserInfo]{}, NO_ERROR
// }

// func (m *mockUserSvc) InviteUser(operUserID string, args model.InviteUserArgs) (SvcRslt[model.InviteUserReply], SvcErr) {
// 	return SvcRslt[model.InviteUserReply]{}, NO_ERROR
// }

// func (m *mockUserSvc) LoginUser(args model.LoginArgs) (SvcRslt[model.LoginReply], SvcErr) {
// 	return SvcRslt[model.LoginReply]{}, NO_ERROR
// }
// func (m *mockUserSvc) UpdateUserInfo(args model.UpdateUserArgs) SvcErr { return NO_ERROR }
// func (m *mockUserSvc) GetUserInfos(opt model.RetrieveUserOpt) (SvcRslt[[]model.UserInfo], SvcErr) {
// 	return SvcRslt[[]model.UserInfo]{}, NO_ERROR
// }

// func TestCreateAndRetrieveWorksets(t *testing.T) {
// 	r := repo.NewMockWorksetRepo()
// 	ws := NewWorksetSvc(r)
// 	// mock user service instance
// 	mu := &mockUserSvc{}

// 	// create two worksets
// 	a1 := &model.CreateWorksetArgs{Name: "First", Description: nil}
// 	a2 := &model.CreateWorksetArgs{Name: "Second", Description: nil}

// 	res1, err1 := ws.CreateWorkset(mu, "u1", a1)
// 	if err1 != NO_ERROR {
// 		t.Fatalf("CreateWorkset a1 failed: %v", err1)
// 	}
// 	if res1.Data == nil || res1.Data.ID == "" {
// 		t.Fatalf("CreateWorkset a1 returned empty id")
// 	}

// 	res2, err2 := ws.CreateWorkset(mu, "u2", a2)
// 	if err2 != NO_ERROR {
// 		t.Fatalf("CreateWorkset a2 failed: %v", err2)
// 	}
// 	if res2.Data == nil || res2.Data.ID == "" {
// 		t.Fatalf("CreateWorkset a2 returned empty id")
// 	}

// 	rslt, svcErr := ws.RetrieveWorksets(0, 0)
// 	if svcErr != NO_ERROR {
// 		t.Fatalf("RetrieveWorksets returned error: %v", svcErr)
// 	}
// 	if rslt.Data == nil || len(*rslt.Data) != 2 {
// 		t.Fatalf("expected 2 worksets, got %v", rslt.Data)
// 	}
// }

// func TestGetAndUpdateWorkset(t *testing.T) {
// 	r := repo.NewMockWorksetRepo()
// 	ws := NewWorksetSvc(r)

// 	// create
// 	a := &model.CreateWorksetArgs{Name: "WU", Description: nil}
// 	mu2 := &mockUserSvc{}
// 	res, err := ws.CreateWorkset(mu2, "creator", a)
// 	if err != NO_ERROR {
// 		t.Fatalf("CreateWorkset failed: %v", err)
// 	}
// 	id := res.Data.ID

// 	g, getErr := ws.GetWorksetByID(id)
// 	if getErr != NO_ERROR {
// 		t.Fatalf("GetWorksetByID failed: %v", getErr)
// 	}
// 	if g.Data.ID != id {
// 		t.Fatalf("unexpected workset id: %v", g.Data.ID)
// 	}

// 	// update description
// 	newDesc := "updated"
// 	updErr := ws.UpdateWorksetByID(&model.UpdateWorksetArgs{ID: id, Description: &newDesc})
// 	if updErr != NO_ERROR {
// 		t.Fatalf("UpdateWorksetByID failed: %v", updErr)
// 	}

// 	g2, getErr2 := ws.GetWorksetByID(id)
// 	if getErr2 != NO_ERROR {
// 		t.Fatalf("GetWorksetByID after update failed: %v", getErr2)
// 	}
// 	if g2.Data.Description == nil || *g2.Data.Description != newDesc {
// 		t.Fatalf("description not updated: expected %v got %v", newDesc, g2.Data.Description)
// 	}
// }
