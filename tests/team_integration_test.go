package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testTeamName   = "backend"
	testTeamNameQA = "qa"
	testUserID1    = "u1"
	testUserID2    = "u2"
	testUserID3    = "u3"
	testUsername1  = "Alice"
	testUsername2  = "Bob"
	testUsername3  = "Rob"
)

func TestTeamHandler_CreateTeam(t *testing.T) {
	t.Run("successfully creates new team", func(t *testing.T) {
		teamSvc := NewFakeTeamService()
		userSvc := NewFakeUserService()
		prSvc := NewFakePRServiceWithUsers(userSvc)
		router := SetupTestRouter(teamSvc, userSvc, prSvc)

		payload := map[string]interface{}{
			"team_name": testTeamName,
			"members": []interface{}{
				map[string]interface{}{
					"user_id":   testUserID1,
					"username":  testUsername1,
					"is_active": true,
				},
				map[string]interface{}{
					"user_id":   testUserID2,
					"username":  testUsername2,
					"is_active": true,
				},
			},
		}

		body, err := json.Marshal(payload)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/team/add", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.True(t, w.Code == http.StatusCreated || w.Code == http.StatusOK,
			"expected status 201 or 200, got %d", w.Code)
	})

	t.Run("returns error when team already exists", func(t *testing.T) {
		teamSvc := NewFakeTeamService()
		userSvc := NewFakeUserService()
		prSvc := NewFakePRServiceWithUsers(userSvc)
		router := SetupTestRouter(teamSvc, userSvc, prSvc)

		payload := map[string]interface{}{
			"team_name": testTeamName,
			"members": []interface{}{
				map[string]interface{}{
					"user_id":   testUserID1,
					"username":  testUsername1,
					"is_active": true,
				},
			},
		}

		body, err := json.Marshal(payload)
		require.NoError(t, err)

		req := httptest.NewRequest(http.MethodPost, "/team/add", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		w1 := httptest.NewRecorder()
		router.ServeHTTP(w1, req)
		require.True(t, w1.Code == http.StatusCreated || w1.Code == http.StatusOK,
			"first request should succeed")

		w2 := httptest.NewRecorder()
		req2 := httptest.NewRequest(http.MethodPost, "/team/add", bytes.NewReader(body))
		req2.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w2, req2)

		assert.Contains(t, w2.Body.String(), "team already exists",
			"expected error about team already existing, got: %s", w2.Body.String())
	})
}

func TestTeamHandler_GetTeamByName(t *testing.T) {
	t.Run("successfully retrieves existing team", func(t *testing.T) {
		teamSvc := NewFakeTeamService()
		userSvc := NewFakeUserService()
		prSvc := NewFakePRServiceWithUsers(userSvc)

		userSvc.registeredUsers[testUserID3] = MakeTestUser(testUserID3, testUsername3, testTeamNameQA, true)

		router := SetupTestRouter(teamSvc, userSvc, prSvc)

		createPayload := map[string]interface{}{
			"team_name": testTeamNameQA,
			"members": []interface{}{
				map[string]interface{}{
					"user_id":   testUserID3,
					"username":  testUsername3,
					"is_active": true,
				},
			},
		}

		createBody, err := json.Marshal(createPayload)
		require.NoError(t, err)

		createReq := httptest.NewRequest(http.MethodPost, "/team/add", bytes.NewReader(createBody))
		createReq.Header.Set("Content-Type", "application/json")
		createW := httptest.NewRecorder()
		router.ServeHTTP(createW, createReq)
		require.True(t, createW.Code == http.StatusCreated || createW.Code == http.StatusOK,
			"team creation should succeed")

		req := httptest.NewRequest(http.MethodGet, "/team/get?team_name="+testTeamNameQA, nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code,
			"expected status 200 for existing team")
	})

	t.Run("returns 200 with empty team for non-existent team", func(t *testing.T) {
		teamSvc := NewFakeTeamService()
		userSvc := NewFakeUserService()
		prSvc := NewFakePRServiceWithUsers(userSvc)
		router := SetupTestRouter(teamSvc, userSvc, prSvc)

		req := httptest.NewRequest(http.MethodGet, "/team/get?team_name=nonexistent", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code,
			"handler returns 200 even for non-existent team")
	})
}
