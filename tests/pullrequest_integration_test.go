package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testPRID       = "pr-1001"
	testPRID2      = "pr-1002"
	testPRName     = "Add search"
	testAuthorID   = "u1"
	testAuthorName = "Alice"
	testTeamDev    = "dev"
)

func TestPullRequestHandler_CreatePR(t *testing.T) {
	t.Run("successfully creates new pull request", func(t *testing.T) {
		teamSvc := NewFakeTeamService()
		userSvc := NewFakeUserService()
		prSvc := NewFakePRServiceWithUsers(userSvc)

		userSvc.registeredUsers[testAuthorID] = MakeTestUser(testAuthorID, testAuthorName, testTeamDev, true)

		router := SetupTestRouter(teamSvc, userSvc, prSvc)

		payload := map[string]interface{}{
			"pull_request_id":   testPRID,
			"pull_request_name": testPRName,
			"author_id":         testAuthorID,
		}

		body, err := json.Marshal(payload)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/pullRequest/create", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.True(t, w.Code == http.StatusCreated || w.Code == http.StatusOK,
			"expected status 201 or 200, got %d", w.Code)
	})

	t.Run("returns error when PR already exists", func(t *testing.T) {
		teamSvc := NewFakeTeamService()
		userSvc := NewFakeUserService()
		prSvc := NewFakePRServiceWithUsers(userSvc)

		userSvc.registeredUsers[testAuthorID] = MakeTestUser(testAuthorID, testAuthorName, testTeamDev, true)

		router := SetupTestRouter(teamSvc, userSvc, prSvc)

		payload := map[string]interface{}{
			"pull_request_id":   testPRID,
			"pull_request_name": testPRName,
			"author_id":         testAuthorID,
		}

		body, err := json.Marshal(payload)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/pullRequest/create", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		w1 := httptest.NewRecorder()
		router.ServeHTTP(w1, req)
		require.True(t, w1.Code == http.StatusCreated || w1.Code == http.StatusOK,
			"first request should succeed")

		w2 := httptest.NewRecorder()
		req2 := httptest.NewRequest(http.MethodPost, "/pullRequest/create", bytes.NewReader(body))
		req2.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w2, req2)

		bodyStr := strings.ToLower(w2.Body.String())
		assert.Contains(t, bodyStr, "pr already exists",
			"expected error about PR already existing, got: %s", w2.Body.String())
	})

	t.Run("returns error when author does not exist", func(t *testing.T) {
		teamSvc := NewFakeTeamService()
		userSvc := NewFakeUserService()
		prSvc := NewFakePRServiceWithUsers(userSvc)
		router := SetupTestRouter(teamSvc, userSvc, prSvc)

		payload := map[string]interface{}{
			"pull_request_id":   testPRID2,
			"pull_request_name": "Wrong user PR",
			"author_id":         "nonexistent_user",
		}

		body, err := json.Marshal(payload)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/pullRequest/create", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		bodyStr := strings.ToLower(w.Body.String())
		assert.Contains(t, bodyStr, "pr not found",
			"expected error about PR/user not found, got: %s", w.Body.String())
	})
}

func TestPullRequestHandler_MergePR(t *testing.T) {
	t.Run("successfully merges existing pull request", func(t *testing.T) {
		teamSvc := NewFakeTeamService()
		userSvc := NewFakeUserService()
		prSvc := NewFakePRServiceWithUsers(userSvc)

		prSvc.createdPRs[testPRID] = MakeTestPR(testPRID, testPRName, testAuthorID, true)

		router := SetupTestRouter(teamSvc, userSvc, prSvc)

		payload := map[string]interface{}{
			"pull_request_id": testPRID,
		}

		body, err := json.Marshal(payload)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/pullRequest/merge", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusCreated,
			"expected status 200 or 201, got %d", w.Code)
	})

	t.Run("returns error when PR does not exist", func(t *testing.T) {
		teamSvc := NewFakeTeamService()
		userSvc := NewFakeUserService()
		prSvc := NewFakePRServiceWithUsers(userSvc)
		router := SetupTestRouter(teamSvc, userSvc, prSvc)

		payload := map[string]interface{}{
			"pull_request_id": "nonexistent_pr",
		}

		body, err := json.Marshal(payload)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/pullRequest/merge", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		bodyStr := strings.ToLower(w.Body.String())
		assert.Contains(t, bodyStr, "pr not found",
			"expected error about PR not found, got: %s", w.Body.String())
	})
}

func TestPullRequestHandler_ReassignPR(t *testing.T) {
	t.Run("successfully reassigns reviewer", func(t *testing.T) {
		teamSvc := NewFakeTeamService()
		userSvc := NewFakeUserService()
		prSvc := NewFakePRServiceWithUsers(userSvc)

		prSvc.createdPRs[testPRID] = MakeTestPR(testPRID, testPRName, testAuthorID, true)

		router := SetupTestRouter(teamSvc, userSvc, prSvc)

		payload := map[string]interface{}{
			"pull_request_id": testPRID,
			"old_user_id":     "rev1",
		}

		body, err := json.Marshal(payload)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/pullRequest/reassign", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusCreated,
			"expected status 200 or 201, got %d", w.Code)
	})

	t.Run("returns error when PR does not exist", func(t *testing.T) {
		teamSvc := NewFakeTeamService()
		userSvc := NewFakeUserService()
		prSvc := NewFakePRServiceWithUsers(userSvc)
		router := SetupTestRouter(teamSvc, userSvc, prSvc)

		payload := map[string]interface{}{
			"pull_request_id": "nonexistent_pr",
			"old_user_id":     "rev1",
		}

		body, err := json.Marshal(payload)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/pullRequest/reassign", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		bodyStr := strings.ToLower(w.Body.String())
		assert.Contains(t, bodyStr, "pr not found",
			"expected error about PR not found, got: %s", w.Body.String())
	})

	t.Run("returns error when PR has no reviewers", func(t *testing.T) {
		teamSvc := NewFakeTeamService()
		userSvc := NewFakeUserService()
		prSvc := NewFakePRServiceWithUsers(userSvc)

		prSvc.createdPRs[testPRID] = MakeTestPR(testPRID, testPRName, testAuthorID, false)

		router := SetupTestRouter(teamSvc, userSvc, prSvc)

		payload := map[string]interface{}{
			"pull_request_id": testPRID,
			"old_user_id":     "nonexistent_reviewer",
		}

		body, err := json.Marshal(payload)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/pullRequest/reassign", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		bodyStr := strings.ToLower(w.Body.String())
		assert.Contains(t, bodyStr, "no reviewers",
			"expected error about no reviewers, got: %s", w.Body.String())
	})
}
