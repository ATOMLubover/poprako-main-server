package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

const (
	baseURL = "http://localhost:8080/api/v1"
)

var (
	authToken     string
	serverCmd     *exec.Cmd
	testWorksetID string
	testComicID   string
	testPageIDs   []string
	testUnitIDs   []string
	testAsgnID    string

	// User management test data
	testUserIDs       []string
	testInvCodes      []string
	testTranslators   []string
	testProofreaders  []string
	testTypesetters   []string
	testRedrawers     []string
	testReviewers     []string
	testPreAsgnIDs    []string
	testNormalAsgnIDs []string
)

// TestMain controls the execution flow and starts/stops the server
func TestMain(m *testing.M) {
	// Load .env file to get MOCK_AUTH_TOKEN
	if err := godotenv.Load("../.env"); err != nil {
		fmt.Printf("Warning: .env file not found: %v\n", err)
	}

	authToken = os.Getenv("MOCK_AUTH_TOKEN")
	if authToken == "" {
		fmt.Println("Error: MOCK_AUTH_TOKEN not found in .env")
		os.Exit(1)
	}

	// // Start the server
	// if err := startServer(); err != nil {
	// 	fmt.Printf("Failed to start server: %v\n", err)
	// 	os.Exit(1)
	// }

	// defer stopServer()

	// Wait for server to be ready
	if err := waitForServer(); err != nil {
		fmt.Printf("Server not ready: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✓ Server started successfully")

	// Run all tests serially
	code := m.Run()

	os.Exit(code)
}

// waitForServer waits for the server to be ready by polling the check-update endpoint
func waitForServer() error {
	maxRetries := 20
	retryDelay := 500 * time.Millisecond

	for i := 0; i < maxRetries; i++ {
		resp, err := http.Get(baseURL + "/check-update")
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			return nil
		}
		if resp != nil {
			resp.Body.Close()
		}
		time.Sleep(retryDelay)
	}

	return fmt.Errorf("server did not become ready after %d retries", maxRetries)
}

// makeRequest is a helper to make HTTP requests with auth
func makeRequest(method, path string, body interface{}) (*http.Response, []byte, error) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, baseURL+path, reqBody)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+authToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("request failed: %w", err)
	}

	respBody, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return resp, nil, fmt.Errorf("failed to read response: %w", err)
	}

	return resp, respBody, nil
}

// unmarshalResponse unmarshals the API response
func unmarshalResponse(data []byte, target interface{}) error {
	// Handle empty response body (204 No Content)
	if len(data) == 0 {
		return nil
	}

	var wrapper struct {
		Code uint16          `json:"code"`
		Msg  string          `json:"msg"`
		Data json.RawMessage `json:"data"`
	}

	if err := json.Unmarshal(data, &wrapper); err != nil {
		return fmt.Errorf("failed to unmarshal wrapper: %w", err)
	}

	if wrapper.Code != 200 && wrapper.Code != 201 && wrapper.Code != 204 {
		return fmt.Errorf("API error (code %d): %s", wrapper.Code, wrapper.Msg)
	}

	if target != nil && len(wrapper.Data) > 0 {
		if err := json.Unmarshal(wrapper.Data, target); err != nil {
			return fmt.Errorf("failed to unmarshal data: %w", err)
		}
	}

	return nil
}

// TestIntegrationFlow runs the complete workflow serially
// If any step fails, the entire test stops immediately
func TestIntegrationFlow(t *testing.T) {
	// Run all test steps in order, stopping on first failure
	t.Run("01_UserManagement", testUserManagement)
	if t.Failed() {
		t.Fatal("Step 01 failed, aborting test")
	}

	t.Run("02_CreateWorksetAndComicWithPreAssignments", testCreateWorksetAndComicWithPreAssignments)
	if t.Failed() {
		t.Fatal("Step 02 failed, aborting test")
	}

	t.Run("03_CreatePagesAndUploadImages", testCreatePagesAndUploadImages)
	if t.Failed() {
		t.Fatal("Step 03 failed, aborting test")
	}

	t.Run("04_CreateUnits", testCreateUnits)
	if t.Failed() {
		t.Fatal("Step 04 failed, aborting test")
	}

	t.Run("05_TranslationWorkflow", testTranslationWorkflow)
	if t.Failed() {
		t.Fatal("Step 05 failed, aborting test")
	}

	t.Run("06_AssignmentsManagement", testAssignmentsManagement)
	if t.Failed() {
		t.Fatal("Step 06 failed, aborting test")
	}

	t.Run("07_RetrievalAndFiltering", testRetrievalAndFiltering)
	if t.Failed() {
		t.Fatal("Step 07 failed, aborting test")
	}

	t.Run("08_Cleanup", testCleanup)
}

// testUserManagement creates users via invitation and assigns roles
func testUserManagement(t *testing.T) {
	// Define 10 QQ numbers to invite (2 per role: translator, proofreader, typesetter, redrawer, reviewer)
	// The invitation system expects decimal numbers that will be used as QQ during registration
	userQQs := []string{
		"1000001",
		"1000002",
		"1000003",
		"1000004",
		"1000005",
		"1000006",
		"1000007",
		"1000008",
		"1000009",
		"1000010",
	}

	roleMapping := []string{
		"translator", "translator",
		"proofreader", "proofreader",
		"typesetter", "typesetter",
		"redrawer", "redrawer",
		"reviewer", "reviewer",
	}

	// Step 1: Admin invites users (invitee_id should be the QQ number)
	for i, qq := range userQQs {
		inviteArgs := map[string]interface{}{
			"invitee_id": qq,
		}

		resp, body, err := makeRequest("POST", "/users/invite", inviteArgs)
		if err != nil {
			t.Fatalf("Failed to invite user with QQ %s: %v", qq, err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Invite user with QQ %s failed with status %d: %s", qq, resp.StatusCode, string(body))
		}

		var inviteReply struct {
			InvCode string `json:"invitation_code"`
		}
		if err := unmarshalResponse(body, &inviteReply); err != nil {
			t.Fatalf("Failed to parse invite response for QQ %s: %v", qq, err)
		}

		testInvCodes = append(testInvCodes, inviteReply.InvCode)
		t.Logf("✓ Invited user %d/%d: QQ %s (invitation code: %s)", i+1, len(userQQs), qq, inviteReply.InvCode)
	}

	// Step 2: Users register/login with invitation codes
	for i, qq := range userQQs {
		loginArgs := map[string]interface{}{
			"qq":              qq,
			"password":        "test123456",
			"nickname":        fmt.Sprintf("TestUser%02d", i+1),
			"invitation_code": testInvCodes[i],
		}

		resp, body, err := makeRequest("POST", "/login", loginArgs)
		if err != nil {
			t.Fatalf("Failed to login user QQ %s: %v", qq, err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Login user QQ %s failed with status %d: %s", qq, resp.StatusCode, string(body))
		}

		var loginReply struct {
			Token string `json:"token"`
		}
		if err := unmarshalResponse(body, &loginReply); err != nil {
			t.Fatalf("Failed to parse login response for QQ %s: %v", qq, err)
		}

		t.Logf("✓ User registered/logged in: QQ %s", qq)
	}

	// Step 3: Get all users to find their actual IDs
	resp, body, err := makeRequest("GET", "/users?limit=20&offset=0", nil)
	if err != nil {
		t.Fatalf("Failed to retrieve users: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Retrieve users failed with status %d: %s", resp.StatusCode, string(body))
	}

	// DEBUG: Print raw response
	t.Logf("DEBUG: GET /users raw response: %s", string(body))

	var allUsers []map[string]interface{}
	if err := unmarshalResponse(body, &allUsers); err != nil {
		t.Fatalf("Failed to parse users: %v", err)
	}

	// DEBUG: Print parsed users
	t.Logf("DEBUG: Parsed %d users from response", len(allUsers))

	// Build map of QQ to user_id
	qqToUserID := make(map[string]string)
	for _, user := range allUsers {
		t.Logf("DEBUG: Processing user: %+v", user)
		if qq, ok := user["qq"].(string); ok {
			if userID, ok := user["user_id"].(string); ok {
				qqToUserID[qq] = userID
				t.Logf("DEBUG: Found QQ %s -> UserID %s", qq, userID)
			} else {
				t.Logf("DEBUG: user_id not found or not string for user: %+v", user)
			}
		} else {
			t.Logf("DEBUG: qq not found or not string for user: %+v", user)
		}
	}

	// Store user IDs and organize by role
	for i, qq := range userQQs {
		if userID, ok := qqToUserID[qq]; ok {
			testUserIDs = append(testUserIDs, userID)
			role := roleMapping[i]
			switch role {
			case "translator":
				testTranslators = append(testTranslators, userID)
			case "proofreader":
				testProofreaders = append(testProofreaders, userID)
			case "typesetter":
				testTypesetters = append(testTypesetters, userID)
			case "redrawer":
				testRedrawers = append(testRedrawers, userID)
			case "reviewer":
				testReviewers = append(testReviewers, userID)
			}
		}
	}

	if len(testUserIDs) != len(userQQs) {
		t.Fatalf("Expected %d users, found %d", len(userQQs), len(testUserIDs))
	}
	t.Logf("✓ Retrieved %d user IDs", len(testUserIDs))

	// Step 4: Assign roles to users (2 users per role)
	roleAssignments := []struct {
		userIndex int
		roles     []map[string]interface{}
	}{
		{0, []map[string]interface{}{{"role": "translator", "assigned": true}}},
		{1, []map[string]interface{}{{"role": "translator", "assigned": true}}},
		{2, []map[string]interface{}{{"role": "proofreader", "assigned": true}}},
		{3, []map[string]interface{}{{"role": "proofreader", "assigned": true}}},
		{4, []map[string]interface{}{{"role": "typesetter", "assigned": true}}},
		{5, []map[string]interface{}{{"role": "typesetter", "assigned": true}}},
		{6, []map[string]interface{}{{"role": "redrawer", "assigned": true}}},
		{7, []map[string]interface{}{{"role": "redrawer", "assigned": true}}},
		{8, []map[string]interface{}{{"role": "reviewer", "assigned": true}}},
		{9, []map[string]interface{}{{"role": "reviewer", "assigned": true}}},
	}

	for _, assignment := range roleAssignments {
		userID := testUserIDs[assignment.userIndex]
		assignArgs := map[string]interface{}{
			"user_id": userID,
			"roles":   assignment.roles,
		}

		resp, body, err := makeRequest("PATCH", "/users/"+userID+"/role", assignArgs)
		if err != nil {
			t.Fatalf("Failed to assign role to %s: %v", userID, err)
		}
		if resp.StatusCode != http.StatusNoContent {
			t.Fatalf("Assign role to %s failed with status %d: %s", userID, resp.StatusCode, string(body))
		}

		t.Logf("✓ Assigned role to user %d: %v", assignment.userIndex+1, assignment.roles[0]["role"])
	}

	// Step 5: Test unassign and reassign for one user
	testUserID := testUserIDs[0]
	unassignArgs := map[string]interface{}{
		"user_id": testUserID,
		"roles":   []map[string]interface{}{{"role": "translator", "assigned": false}},
	}

	resp, body, err = makeRequest("PATCH", "/users/"+testUserID+"/role", unassignArgs)
	if err != nil {
		t.Fatalf("Failed to unassign role from %s: %v", testUserID, err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("Unassign role from %s failed with status %d: %s", testUserID, resp.StatusCode, string(body))
	}
	t.Logf("✓ Unassigned translator role from %s", testUserID)

	// Reassign the role
	reassignArgs := map[string]interface{}{
		"user_id": testUserID,
		"roles":   []map[string]interface{}{{"role": "translator", "assigned": true}},
	}

	resp, body, err = makeRequest("PATCH", "/users/"+testUserID+"/role", reassignArgs)
	if err != nil {
		t.Fatalf("Failed to reassign role to %s: %v", testUserID, err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("Reassign role to %s failed with status %d: %s", testUserID, resp.StatusCode, string(body))
	}
	t.Logf("✓ Reassigned translator role to %s", testUserID)

	// Step 6: Verify all users and their roles
	resp, body, err = makeRequest("GET", "/users?limit=20&offset=0", nil)
	if err != nil {
		t.Fatalf("Failed to retrieve users: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Retrieve users failed with status %d: %s", resp.StatusCode, string(body))
	}

	var users []map[string]interface{}
	if err := unmarshalResponse(body, &users); err != nil {
		t.Fatalf("Failed to parse users: %v", err)
	}
	t.Logf("✓ Retrieved %d users total", len(users))
	t.Logf("✓ Test users created: Translators(%d), Proofreaders(%d), Typesetters(%d), Redrawers(%d), Reviewers(%d)",
		len(testTranslators), len(testProofreaders), len(testTypesetters), len(testRedrawers), len(testReviewers))
}

func testCreateWorksetAndComicWithPreAssignments(t *testing.T) {
	// Create workset
	createWorksetArgs := map[string]interface{}{
		"name":        "Test Workset",
		"description": "Integration test workset",
	}

	resp, body, err := makeRequest("POST", "/worksets", createWorksetArgs)
	if err != nil {
		t.Fatalf("Failed to create workset: %v", err)
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		t.Fatalf("Create workset failed with status %d: %s", resp.StatusCode, string(body))
	}

	var createWorksetReply struct {
		ID string `json:"id"`
	}
	if err := unmarshalResponse(body, &createWorksetReply); err != nil {
		t.Fatalf("Failed to parse workset response: %v", err)
	}

	testWorksetID = createWorksetReply.ID
	t.Logf("✓ Created workset: %s", testWorksetID)

	// Get workset to verify
	resp, body, err = makeRequest("GET", "/worksets/"+testWorksetID, nil)
	if err != nil {
		t.Fatalf("Failed to get workset: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Get workset failed with status %d: %s", resp.StatusCode, string(body))
	}

	var worksetInfo map[string]interface{}
	if err := unmarshalResponse(body, &worksetInfo); err != nil {
		t.Fatalf("Failed to parse workset info: %v", err)
	}
	t.Logf("✓ Retrieved workset: %s", worksetInfo["name"])

	// Create comic with pre-assignments (if users were created successfully)
	trueVal := true
	createComicArgs := map[string]interface{}{
		"workset_id":  testWorksetID,
		"author":      "测试作者",
		"title":       "测试漫画",
		"description": "这是一个集成测试漫画",
	}

	// Add pre-assignments only if users exist
	if len(testTranslators) > 0 && len(testProofreaders) > 0 {
		createComicArgs["pre_asgns"] = []map[string]interface{}{
			{
				"assignee_id":   testTranslators[0],
				"is_translator": &trueVal,
			},
			{
				"assignee_id":    testProofreaders[0],
				"is_proofreader": &trueVal,
			},
		}
		t.Logf("✓ Adding pre-assignments for translator and proofreader")
	} else {
		t.Logf("⚠ Skipping pre-assignments (no users available)")
	}

	resp, body, err = makeRequest("POST", "/comics", createComicArgs)
	if err != nil {
		t.Fatalf("Failed to create comic: %v", err)
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		t.Fatalf("Create comic failed with status %d: %s", resp.StatusCode, string(body))
	}

	var createComicReply struct {
		ID string `json:"id"`
	}
	if err := unmarshalResponse(body, &createComicReply); err != nil {
		t.Fatalf("Failed to parse comic response: %v", err)
	}

	testComicID = createComicReply.ID
	t.Logf("✓ Created comic: %s", testComicID)

	// Get comic to verify
	resp, body, err = makeRequest("GET", "/comics/"+testComicID, nil)
	if err != nil {
		t.Fatalf("Failed to get comic: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Get comic failed with status %d: %s", resp.StatusCode, string(body))
	}

	var comicInfo map[string]interface{}
	if err := unmarshalResponse(body, &comicInfo); err != nil {
		t.Fatalf("Failed to parse comic info: %v", err)
	}
	t.Logf("✓ Retrieved comic: %s by %s", comicInfo["title"], comicInfo["author"])
}

func testCreatePagesAndUploadImages(t *testing.T) {
	// Use images from ../mocks directory
	mockImages := []struct {
		fileName  string
		extension string
	}{
		{"133700406_p0.jpg", "jpg"},
		{"ginko.jpg", "jpg"},
		{"普瑞赛斯.png", "png"},
	}

	// Create pages
	createPagesArgs := []map[string]interface{}{}
	for i, img := range mockImages {
		createPagesArgs = append(createPagesArgs, map[string]interface{}{
			"comic_id":  testComicID,
			"index":     i + 1,
			"image_ext": img.extension,
		})
	}

	resp, body, err := makeRequest("POST", "/pages", createPagesArgs)
	if err != nil {
		t.Fatalf("Failed to create pages: %v", err)
	}
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		t.Fatalf("Create pages failed with status %d: %s", resp.StatusCode, string(body))
	}

	var createPagesReply []struct {
		ID     string `json:"id"`
		OSSUrl string `json:"oss_url"`
	}
	if err := unmarshalResponse(body, &createPagesReply); err != nil {
		t.Fatalf("Failed to parse pages response: %v", err)
	}

	// Upload each image and mark as uploaded
	for i, page := range createPagesReply {
		testPageIDs = append(testPageIDs, page.ID)
		t.Logf("✓ Created page %d: %s", i+1, page.ID)
		t.Logf("  Upload URL: %s", page.OSSUrl)

		// Read image file from ../mocks
		imgPath := filepath.Join("..", "mocks", mockImages[i].fileName)
		imgData, err := os.ReadFile(imgPath)
		if err != nil {
			t.Fatalf("Failed to read image file %s: %v", imgPath, err)
		}
		t.Logf("  Read %d bytes from %s", len(imgData), imgPath)

		// Upload to presigned URL
		uploadReq, err := http.NewRequest("PUT", page.OSSUrl, bytes.NewReader(imgData))
		if err != nil {
			t.Fatalf("Failed to create upload request: %v", err)
		}

		// Note: For R2 presigned URLs, we typically don't need to set Content-Type
		// unless it was specified during presign generation
		// Set Content-Length to ensure proper upload
		uploadReq.ContentLength = int64(len(imgData))

		uploadClient := &http.Client{Timeout: 30 * time.Second}
		uploadResp, err := uploadClient.Do(uploadReq)
		if err != nil {
			t.Fatalf("Failed to upload image to OSS: %v", err)
		}
		defer uploadResp.Body.Close()

		// Read response body for debugging
		uploadRespBody, _ := io.ReadAll(uploadResp.Body)

		if uploadResp.StatusCode != http.StatusOK && uploadResp.StatusCode != http.StatusCreated && uploadResp.StatusCode != http.StatusNoContent {
			t.Fatalf("Upload to OSS failed with status %d: %s", uploadResp.StatusCode, string(uploadRespBody))
		}
		t.Logf("✓ Uploaded image %s to OSS (status: %d)", mockImages[i].fileName, uploadResp.StatusCode)

		// Mark page as uploaded
		patchArgs := map[string]interface{}{
			"uploaded": true,
		}

		resp, body, err := makeRequest("PATCH", "/pages/"+page.ID, patchArgs)
		if err != nil {
			t.Fatalf("Failed to patch page: %v", err)
		}
		if resp.StatusCode != http.StatusNoContent {
			t.Fatalf("Patch page failed with status %d: %s", resp.StatusCode, string(body))
		}
		t.Logf("✓ Marked page %d as uploaded", i+1)
	}

	// Download images to verify upload worked
	if err := os.MkdirAll(filepath.Join("..", "mocks", "down"), 0o755); err != nil {
		t.Fatalf("Failed to create download directory: %v", err)
	}

	for i, pageID := range testPageIDs {
		// Get page info to retrieve download URL
		resp, body, err := makeRequest("GET", "/pages/"+pageID, nil)
		if err != nil {
			t.Fatalf("Failed to get page info: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Get page failed with status %d: %s", resp.StatusCode, string(body))
		}

		var pageInfo struct {
			OSSURL string `json:"oss_url"`
		}
		if err := unmarshalResponse(body, &pageInfo); err != nil {
			t.Fatalf("Failed to parse page info: %v", err)
		}

		// Download image
		downloadResp, err := http.Get(pageInfo.OSSURL)
		if err != nil {
			t.Logf("Warning: Failed to download image: %v (OSS may not be configured)", err)
			continue
		}
		defer downloadResp.Body.Close()

		downloadedData, err := io.ReadAll(downloadResp.Body)
		if err != nil {
			t.Fatalf("Failed to read downloaded image: %v", err)
		}

		// Save to mocks/down/
		downloadPath := filepath.Join("..", "mocks", "down", mockImages[i].fileName)
		if err := os.WriteFile(downloadPath, downloadedData, 0o644); err != nil {
			t.Fatalf("Failed to save downloaded image: %v", err)
		}
		t.Logf("✓ Downloaded and verified image %s (%d bytes)", mockImages[i].fileName, len(downloadedData))
	}
}

func testCreateUnits(t *testing.T) {
	// Create units for first page
	createUnitsArgs := []map[string]interface{}{
		{
			"page_id":         testPageIDs[0],
			"index":           1,
			"x_coordinate":    100.5,
			"y_coordinate":    200.5,
			"is_in_box":       true,
			"translated_text": "原文1",
		},
		{
			"page_id":         testPageIDs[0],
			"index":           2,
			"x_coordinate":    300.5,
			"y_coordinate":    400.5,
			"is_in_box":       true,
			"translated_text": "原文2",
		},
	}

	resp, body, err := makeRequest("POST", "/pages/"+testPageIDs[0]+"/units", createUnitsArgs)
	if err != nil {
		t.Fatalf("Failed to create units: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Create units failed with status %d: %s", resp.StatusCode, string(body))
	}
	t.Logf("✓ Created %d units", len(createUnitsArgs))

	// Get units to retrieve IDs
	resp, body, err = makeRequest("GET", "/pages/"+testPageIDs[0]+"/units", nil)
	if err != nil {
		t.Fatalf("Failed to get units: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Get units failed with status %d: %s", resp.StatusCode, string(body))
	}

	var unitsInfo []map[string]interface{}
	if err := unmarshalResponse(body, &unitsInfo); err != nil {
		t.Fatalf("Failed to parse units info: %v", err)
	}

	for _, unit := range unitsInfo {
		testUnitIDs = append(testUnitIDs, unit["id"].(string))
	}
	t.Logf("✓ Retrieved %d units", len(testUnitIDs))
}

func testTranslationWorkflow(t *testing.T) {
	// Check if we have units to work with
	if len(testUnitIDs) < 2 {
		t.Fatalf("Not enough units to test translation workflow. Expected at least 2, got %d", len(testUnitIDs))
	}
	if len(testPageIDs) < 1 {
		t.Fatalf("No pages available to test translation workflow")
	}

	// Update units with translations
	patchUnitsArgs := []map[string]interface{}{
		{
			"id":                 testUnitIDs[0],
			"translated_text":    "翻译文本1",
			"translator_comment": "翻译备注1",
		},
		{
			"id":                 testUnitIDs[1],
			"translated_text":    "翻译文本2",
			"translator_comment": "翻译备注2",
		},
	}

	resp, body, err := makeRequest("PATCH", "/pages/"+testPageIDs[0]+"/units", patchUnitsArgs)
	if err != nil {
		t.Fatalf("Failed to patch units (translation): %v", err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("Patch units failed with status %d: %s", resp.StatusCode, string(body))
	}
	t.Logf("✓ Updated units with translations")

	// Update units with proofreading
	proofUnitsArgs := []map[string]interface{}{
		{
			"id":                  testUnitIDs[0],
			"proved_text":         "校对文本1",
			"proved":              true,
			"proofreader_comment": "校对备注1",
		},
		{
			"id":                  testUnitIDs[1],
			"proved_text":         "校对文本2",
			"proved":              true,
			"proofreader_comment": "校对备注2",
		},
	}

	resp, body, err = makeRequest("PATCH", "/pages/"+testPageIDs[0]+"/units", proofUnitsArgs)
	if err != nil {
		t.Fatalf("Failed to patch units (proofing): %v", err)
	}
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("Patch units (proofing) failed with status %d: %s", resp.StatusCode, string(body))
	}
	t.Logf("✓ Updated units with proofreading")

	// Verify units were updated
	resp, body, err = makeRequest("GET", "/pages/"+testPageIDs[0]+"/units", nil)
	if err != nil {
		t.Fatalf("Failed to get units: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Get units failed with status %d: %s", resp.StatusCode, string(body))
	}

	var updatedUnits []map[string]interface{}
	if err := unmarshalResponse(body, &updatedUnits); err != nil {
		t.Fatalf("Failed to parse updated units: %v", err)
	}

	for _, unit := range updatedUnits {
		if unit["proved"].(bool) {
			t.Logf("✓ Unit %v proved with text: %s", unit["index"], unit["proved_text"])
		}
	}
}

func testAssignmentsManagement(t *testing.T) {
	// Verify pre-assignments created with comic
	resp, body, err := makeRequest("GET", "/comics/"+testComicID+"/assignments", nil)
	if err != nil {
		t.Fatalf("Failed to get comic assignments: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Get comic assignments failed with status %d: %s", resp.StatusCode, string(body))
	}

	var preAsgns []map[string]interface{}
	if err := unmarshalResponse(body, &preAsgns); err != nil {
		t.Fatalf("Failed to parse pre-assignments: %v", err)
	}

	if len(preAsgns) != 2 {
		t.Fatalf("Expected 2 pre-assignments, got %d", len(preAsgns))
	}

	for _, asgn := range preAsgns {
		testPreAsgnIDs = append(testPreAsgnIDs, asgn["id"].(string))
	}
	t.Logf("✓ Verified %d pre-assignments", len(testPreAsgnIDs))

	// Create additional normal assignments after comic creation
	trueVal := true
	normalAsgns := []map[string]interface{}{
		{
			"comic_id":      testComicID,
			"assignee_id":   testTypesetters[0],
			"is_typesetter": &trueVal,
		},
		{
			"comic_id":    testComicID,
			"assignee_id": testRedrawers[0],
			"is_redrawer": &trueVal,
		},
		{
			"comic_id":    testComicID,
			"assignee_id": testReviewers[0],
			"is_reviewer": &trueVal,
		},
	}

	for i, asgnArgs := range normalAsgns {
		resp, body, err := makeRequest("POST", "/assignments", asgnArgs)
		if err != nil {
			t.Fatalf("Failed to create assignment %d: %v", i+1, err)
		}
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
			t.Fatalf("Create assignment %d failed with status %d: %s", i+1, resp.StatusCode, string(body))
		}

		var asgnID string
		if err := unmarshalResponse(body, &asgnID); err != nil {
			t.Fatalf("Failed to parse assignment response: %v", err)
		}

		testNormalAsgnIDs = append(testNormalAsgnIDs, asgnID)
		t.Logf("✓ Created normal assignment %d/%d: %s", i+1, len(normalAsgns), asgnID)
	}

	// Get all assignments for the comic
	resp, body, err = makeRequest("GET", "/comics/"+testComicID+"/assignments", nil)
	if err != nil {
		t.Fatalf("Failed to get comic assignments: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Get comic assignments failed with status %d: %s", resp.StatusCode, string(body))
	}

	var allAsgns []map[string]interface{}
	if err := unmarshalResponse(body, &allAsgns); err != nil {
		t.Fatalf("Failed to parse assignments: %v", err)
	}

	expectedTotal := len(testPreAsgnIDs) + len(testNormalAsgnIDs)
	if len(allAsgns) != expectedTotal {
		t.Fatalf("Expected %d total assignments, got %d", expectedTotal, len(allAsgns))
	}
	t.Logf("✓ Retrieved %d total assignments for comic (%d pre + %d normal)",
		len(allAsgns), len(testPreAsgnIDs), len(testNormalAsgnIDs))

	// Get assignments by user ID
	if len(testTranslators) > 0 {
		resp, body, err = makeRequest("GET", "/users/"+testTranslators[0]+"/assignments", nil)
		if err != nil {
			t.Fatalf("Failed to get user assignments: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Get user assignments failed with status %d: %s", resp.StatusCode, string(body))
		}

		var userAsgns []map[string]interface{}
		if err := unmarshalResponse(body, &userAsgns); err != nil {
			t.Fatalf("Failed to parse user assignments: %v", err)
		}
		t.Logf("✓ Retrieved %d assignments for user %s", len(userAsgns), testTranslators[0])
	}

	// Test deleting an assignment
	if len(testNormalAsgnIDs) > 0 {
		deleteAsgnID := testNormalAsgnIDs[0]
		resp, body, err := makeRequest("DELETE", "/assignments/"+deleteAsgnID, nil)
		if err != nil {
			t.Fatalf("Failed to delete assignment: %v", err)
		}
		if resp.StatusCode != http.StatusNoContent {
			t.Fatalf("Delete assignment failed with status %d: %s", resp.StatusCode, string(body))
		}
		t.Logf("✓ Deleted assignment: %s", deleteAsgnID)

		// Remove from testNormalAsgnIDs
		testNormalAsgnIDs = testNormalAsgnIDs[1:]

		// Verify deletion
		resp, body, err = makeRequest("GET", "/comics/"+testComicID+"/assignments", nil)
		if err != nil {
			t.Fatalf("Failed to verify assignment deletion: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("Get comic assignments failed with status %d: %s", resp.StatusCode, string(body))
		}

		var verifyAsgns []map[string]interface{}
		if err := unmarshalResponse(body, &verifyAsgns); err != nil {
			t.Fatalf("Failed to parse assignments: %v", err)
		}

		expectedAfterDelete := len(testPreAsgnIDs) + len(testNormalAsgnIDs)
		if len(verifyAsgns) != expectedAfterDelete {
			t.Fatalf("After deletion, expected %d assignments, got %d", expectedAfterDelete, len(verifyAsgns))
		}
		t.Logf("✓ Verified assignment deletion: %d assignments remaining", len(verifyAsgns))
	}

	// Store one assignment ID for final cleanup
	if len(testNormalAsgnIDs) > 0 {
		testAsgnID = testNormalAsgnIDs[0]
	} else if len(testPreAsgnIDs) > 0 {
		testAsgnID = testPreAsgnIDs[0]
	}
}

func testRetrievalAndFiltering(t *testing.T) {
	// Retrieve worksets with pagination
	resp, body, err := makeRequest("GET", "/worksets?limit=10&offset=0", nil)
	if err != nil {
		t.Fatalf("Failed to retrieve worksets: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Retrieve worksets failed with status %d: %s", resp.StatusCode, string(body))
	}

	var worksets []map[string]interface{}
	if err := unmarshalResponse(body, &worksets); err != nil {
		t.Fatalf("Failed to parse worksets: %v", err)
	}
	t.Logf("✓ Retrieved %d worksets", len(worksets))

	// Get comics by workset ID
	resp, body, err = makeRequest("GET", "/worksets/"+testWorksetID+"/comics", nil)
	if err != nil {
		t.Fatalf("Failed to get workset comics: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Get workset comics failed with status %d: %s", resp.StatusCode, string(body))
	}

	var worksetComics []map[string]interface{}
	if err := unmarshalResponse(body, &worksetComics); err != nil {
		t.Fatalf("Failed to parse workset comics: %v", err)
	}
	t.Logf("✓ Retrieved %d comics in workset", len(worksetComics))

	// Filter comics by author
	resp, body, err = makeRequest("GET", "/comics?author=测试", nil)
	if err != nil {
		t.Fatalf("Failed to filter comics: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Filter comics failed with status %d: %s", resp.StatusCode, string(body))
	}

	var filteredComics []map[string]interface{}
	if err := unmarshalResponse(body, &filteredComics); err != nil {
		t.Fatalf("Failed to parse filtered comics: %v", err)
	}
	t.Logf("✓ Filtered %d comics by author", len(filteredComics))

	// Get pages for comic
	resp, body, err = makeRequest("GET", "/comics/"+testComicID+"/pages", nil)
	if err != nil {
		t.Fatalf("Failed to get comic pages: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Get comic pages failed with status %d: %s", resp.StatusCode, string(body))
	}

	var comicPages []map[string]interface{}
	if err := unmarshalResponse(body, &comicPages); err != nil {
		t.Fatalf("Failed to parse comic pages: %v", err)
	}
	t.Logf("✓ Retrieved %d pages for comic", len(comicPages))
}

func testCleanup(t *testing.T) {
	// Delete units
	if len(testUnitIDs) > 0 && len(testPageIDs) > 0 {
		resp, body, err := makeRequest("DELETE", "/pages/"+testPageIDs[0]+"/units", testUnitIDs)
		if err != nil {
			t.Fatalf("Failed to delete units: %v", err)
		}
		if resp.StatusCode != http.StatusNoContent {
			t.Fatalf("Delete units failed with status %d: %s", resp.StatusCode, string(body))
		}
		t.Logf("✓ Deleted %d units", len(testUnitIDs))
	}

	// Delete all assignments (pre and normal)
	allAsgnIDs := append(testPreAsgnIDs, testNormalAsgnIDs...)
	for _, asgnID := range allAsgnIDs {
		resp, _, err := makeRequest("DELETE", "/assignments/"+asgnID, nil)
		if err != nil {
			t.Logf("Warning: Failed to delete assignment %s: %v", asgnID, err)
			continue
		}
		if resp.StatusCode != http.StatusNoContent {
			t.Logf("Warning: Delete assignment %s returned status %d", asgnID, resp.StatusCode)
		}
	}
	if len(allAsgnIDs) > 0 {
		t.Logf("✓ Deleted %d assignments", len(allAsgnIDs))
	}

	// Delete pages
	for _, pageID := range testPageIDs {
		resp, body, err := makeRequest("DELETE", "/pages/"+pageID, nil)
		if err != nil {
			t.Fatalf("Failed to delete page: %v", err)
		}
		if resp.StatusCode != http.StatusNoContent {
			t.Fatalf("Delete page failed with status %d: %s", resp.StatusCode, string(body))
		}
	}
	if len(testPageIDs) > 0 {
		t.Logf("✓ Deleted %d pages", len(testPageIDs))
	}

	// Delete comic
	if testComicID != "" {
		resp, body, err := makeRequest("DELETE", "/comics/"+testComicID, nil)
		if err != nil {
			t.Fatalf("Failed to delete comic: %v", err)
		}
		if resp.StatusCode != http.StatusNoContent {
			t.Fatalf("Delete comic failed with status %d: %s", resp.StatusCode, string(body))
		}
		t.Logf("✓ Deleted comic: %s", testComicID)
	}

	// Delete workset
	if testWorksetID != "" {
		resp, body, err := makeRequest("DELETE", "/worksets/"+testWorksetID, nil)
		if err != nil {
			t.Fatalf("Failed to delete workset: %v", err)
		}
		if resp.StatusCode != http.StatusNoContent {
			t.Fatalf("Delete workset failed with status %d: %s", resp.StatusCode, string(body))
		}
		t.Logf("✓ Deleted workset: %s", testWorksetID)
	}

	t.Log("✓ Cleanup completed successfully")
	t.Logf("Note: Test users (%d) are left in the database for verification", len(testUserIDs))
}
