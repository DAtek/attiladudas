package app_fiar

import (
	"attiladudas/backend/ws"
	ws_mocks "attiladudas/backend/ws/mocks"
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DAtek/gotils"
	"github.com/stretchr/testify/assert"
)

func TestFiveInARow(t *testing.T) {
	path := "/ws/five-in-a-row/"
	t.Run("Returns bad request if connection upgrade fails", func(t *testing.T) {
		app := newMockApp()
		app.wsUpgrader.upgrade = func(
			w http.ResponseWriter,
			r *http.Request,
			responseHeader http.Header,
		) (ws.IConn, error) {
			return nil, errors.New("Failed to upgrade")
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(
			"GET",
			path,
			&bytes.Buffer{},
		)
		app.engine.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Room manager is handling the WS connection", func(t *testing.T) {
		app := newMockApp()
		conn := ws_mocks.NewMockChanConn()
		timeout := gotils.NewTimeoutMs(10000)
		go func() { panic(<-timeout.ErrorCh) }()
		defer timeout.Cancel()

		roomManagerWasCalled := false

		app.roomManager.handleConnection = func(connection ws.IConn) {
			roomManagerWasCalled = true
			assert.Equal(t, conn, connection)
		}

		app.wsUpgrader.upgrade = func(
			w http.ResponseWriter,
			r *http.Request,
			responseHeader http.Header,
		) (ws.IConn, error) {
			return conn, nil
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(
			"GET",
			path,
			&bytes.Buffer{},
		)
		app.engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.True(t, roomManagerWasCalled)
	})
}
