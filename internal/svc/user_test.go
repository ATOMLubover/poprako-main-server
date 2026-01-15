package svc

// import (
// 	"os"
// 	"testing"
// 	"time"

// 	"poprako-main-server/internal/jwtcodec"
// 	"poprako-main-server/internal/model"
// 	"poprako-main-server/internal/model/po"
// 	"poprako-main-server/internal/repo"
// )

// func TestInviteAndRegisterFlow(t *testing.T) {
// 	r := repo.NewMockUserRepo()

// 	os.Setenv("JWT_SECRET_KEY", "testkey")
// 	jc := jwtcodec.NewJWTCodec(3600)

// 	usIfc := NewUserSvc(r, jc)
// 	us, ok := usIfc.(*userSvc)
// 	if !ok {
// 		t.Fatalf("failed to get concrete userSvc")
// 	}

// 	// prepare admin user
// 	adminID := "admin1"
// 	if err := r.CreateUser(nil, &po.NewUser{ID: adminID, QQ: "9999", Nickname: "admin", PasswordHash: "x"}); err != nil {
// 		t.Fatalf("failed to create admin: %v", err)
// 	}
// 	isAdmin := true
// 	if err := r.UpdateUserByID(nil, &po.PatchUser{ID: adminID, IsAdmin: &isAdmin}); err != nil {
// 		t.Fatalf("failed to set admin flag: %v", err)
// 	}

// 	inviteeQQ := "1234"
// 	rslt, svcErr := us.InviteUser(adminID, model.InviteUserArgs{InviteeID: inviteeQQ})
// 	if svcErr != NO_ERROR {
// 		t.Fatalf("InviteUser returned unexpected error: %v", svcErr)
// 	}
// 	if rslt.Data == nil || rslt.Data.InvCode == "" {
// 		t.Fatalf("expected invitation code, got none")
// 	}

// 	invCode := rslt.Data.InvCode

// 	// Now login/register with invitation code using exported interface
// 	loginArgs := model.LoginArgs{QQ: inviteeQQ, Password: "passwd", Nickname: "nick", InvCode: invCode}
// 	loginRslt, loginErr := us.LoginUser(loginArgs)
// 	if loginErr != NO_ERROR {
// 		t.Fatalf("LoginUser (register) returned error: %v", loginErr)
// 	}
// 	if loginRslt.Data == nil || loginRslt.Data.Token == "" {
// 		t.Fatalf("expected token on successful registration, got none")
// 	}
// }

// func TestInvitePermissionDenied(t *testing.T) {
// 	r := repo.NewMockUserRepo()
// 	os.Setenv("JWT_SECRET_KEY", "testkey")
// 	jc := jwtcodec.NewJWTCodec(3600)
// 	usIfc := NewUserSvc(r, jc)
// 	us, ok := usIfc.(*userSvc)
// 	if !ok {
// 		t.Fatalf("failed to get concrete userSvc")
// 	}

// 	// create a non-admin user
// 	userID := "u1"
// 	if err := r.CreateUser(nil, &po.NewUser{ID: userID, QQ: "2000", Nickname: "user", PasswordHash: "x"}); err != nil {
// 		t.Fatalf("failed to create user: %v", err)
// 	}

// 	_, svcErr := us.InviteUser(userID, model.InviteUserArgs{InviteeID: "3000"})
// 	if svcErr != PERMISSION_DENIED {
// 		t.Fatalf("expected PERMISSION_DENIED, got %v", svcErr)
// 	}
// }

// func TestLoginExistingSuccessAndDecodeJWT(t *testing.T) {
// 	r := repo.NewMockUserRepo()
// 	os.Setenv("JWT_SECRET_KEY", "testkey")
// 	jc := jwtcodec.NewJWTCodec(3600)

// 	usIfc := NewUserSvc(r, jc)
// 	us, ok := usIfc.(*userSvc)
// 	if !ok {
// 		t.Fatalf("failed to get concrete userSvc")
// 	}

// 	// create user with hashed password
// 	uid := "exist1"
// 	pwd := "s3cr3t"
// 	hashed, err := us.hashPwd(pwd)
// 	if err != nil {
// 		t.Fatalf("hashPwd failed: %v", err)
// 	}
// 	if err := r.CreateUser(nil, &po.NewUser{ID: uid, QQ: "4000", Nickname: "euser", PasswordHash: hashed}); err != nil {
// 		t.Fatalf("CreateUser failed: %v", err)
// 	}

// 	// login
// 	loginRslt, loginErr := us.LoginUser(model.LoginArgs{QQ: "4000", Password: pwd})
// 	if loginErr != NO_ERROR {
// 		t.Fatalf("LoginUser returned error: %v", loginErr)
// 	}
// 	if loginRslt.Data == nil || loginRslt.Data.Token == "" {
// 		t.Fatalf("expected token on login, got none")
// 	}

// 	// decode token
// 	claims, err := jc.Decode(loginRslt.Data.Token)
// 	if err != nil {
// 		t.Fatalf("failed to decode token: %v", err)
// 	}
// 	if claims.UserID != uid {
// 		t.Fatalf("token contains unexpected userid: %s", claims.UserID)
// 	}
// }

// func TestLoginPasswordMismatch(t *testing.T) {
// 	r := repo.NewMockUserRepo()
// 	os.Setenv("JWT_SECRET_KEY", "testkey")
// 	jc := jwtcodec.NewJWTCodec(3600)

// 	usIfc := NewUserSvc(r, jc)
// 	us, ok := usIfc.(*userSvc)
// 	if !ok {
// 		t.Fatalf("failed to get concrete userSvc")
// 	}

// 	uid := "exist2"
// 	pwd := "rightpwd"
// 	hashed, err := us.hashPwd(pwd)
// 	if err != nil {
// 		t.Fatalf("hashPwd failed: %v", err)
// 	}
// 	if err := r.CreateUser(nil, &po.NewUser{ID: uid, QQ: "5000", Nickname: "euser2", PasswordHash: hashed}); err != nil {
// 		t.Fatalf("CreateUser failed: %v", err)
// 	}

// 	// attempt login with wrong password
// 	_, loginErr := us.LoginUser(model.LoginArgs{QQ: "5000", Password: "wrongpwd"})
// 	if loginErr != PWD_MISMATCH {
// 		t.Fatalf("expected PWD_MISMATCH, got %v", loginErr)
// 	}
// }

// func TestGetAndUpdateUserInfo(t *testing.T) {
// 	r := repo.NewMockUserRepo()
// 	os.Setenv("JWT_SECRET_KEY", "testkey")
// 	jc := jwtcodec.NewJWTCodec(3600)
// 	usIfc := NewUserSvc(r, jc)
// 	us, ok := usIfc.(*userSvc)
// 	if !ok {
// 		t.Fatalf("failed to get concrete userSvc")
// 	}

// 	uid := "upd1"
// 	if err := r.CreateUser(nil, &po.NewUser{ID: uid, QQ: "6000", Nickname: "oldnick", PasswordHash: "x"}); err != nil {
// 		t.Fatalf("CreateUser failed: %v", err)
// 	}

// 	// get by id
// 	infoRslt, err := us.GetUserInfoByID(uid)
// 	if err != NO_ERROR {
// 		t.Fatalf("GetUserInfoByID error: %v", err)
// 	}
// 	if infoRslt.Data == nil || infoRslt.Data.Nickname != "oldnick" {
// 		t.Fatalf("unexpected user info: %+v", infoRslt.Data)
// 	}

// 	// update nickname and qq
// 	newQQ := "6001"
// 	newNick := "newnick"
// 	if updErr := us.UpdateUserInfo(model.UpdateUserArgs{UserID: uid, QQ: &newQQ, Nickname: &newNick}); updErr != NO_ERROR {
// 		t.Fatalf("UpdateUserInfo failed: %v", updErr)
// 	}

// 	// allow a small delay if timestamps etc used
// 	time.Sleep(10 * time.Millisecond)

// 	// check updated by QQ
// 	infoRslt2, err2 := us.GetUserInfoByQQ(newQQ)
// 	if err2 != NO_ERROR {
// 		t.Fatalf("GetUserInfoByQQ error: %v", err2)
// 	}
// 	if infoRslt2.Data == nil || infoRslt2.Data.Nickname != newNick {
// 		t.Fatalf("unexpected updated user info: %+v", infoRslt2.Data)
// 	}
// }
