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

const testTeamBackend = "backend"

func TestUserHandler_SetActive(t *testing.T) {
	t.Run("successfully updates user active status", func(t *testing.T) {
		teamSvc := NewFakeTeamService()
		userSvc := NewFakeUserService()
		prSvc := NewFakePRServiceWithUsers(userSvc)

		userSvc.registeredUsers[testUserID2] = MakeTestUser(testUserID2, testUsername2, testTeamBackend, true)

		router := SetupTestRouter(teamSvc, userSvc, prSvc)

		payload := map[string]interface{}{
			"user_id":   testUserID2,
			"is_active": false,
		}

		body, err := json.Marshal(payload)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/users/setIsActive", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusCreated,
			"expected status 200 or 201, got %d", w.Code)
	})

	t.Run("returns error when user does not exist", func(t *testing.T) {
		teamSvc := NewFakeTeamService()
		userSvc := NewFakeUserService()
		prSvc := NewFakePRServiceWithUsers(userSvc)
		router := SetupTestRouter(teamSvc, userSvc, prSvc)

		payload := map[string]interface{}{
			"user_id":   "nonexistent_user",
			"is_active": false,
		}

		body, err := json.Marshal(payload)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/users/setIsActive", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		bodyStr := strings.ToLower(w.Body.String())
		assert.Contains(t, bodyStr, "user not found",
			"expected error message about user not found, got: %s", w.Body.String())
	})
}
