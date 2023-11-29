package fiar

import (
	"api/components/room_manager"
	"bytes"
	"fibertools"
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/DAtek/gotils"
	"github.com/stretchr/testify/assert"
)

func TestFiveInARow(t *testing.T) {
	t.Run("Room manager is handling the WS connection", func(t *testing.T) {
		timeout := gotils.NewTimeoutMs(10000)
		defer timeout.Cancel()
		go func() { panic(<-timeout.ErrorCh) }()

		called := false
		roomManager := &room_manager.MockRoomManager{
			HandleConnection_: func(conn room_manager.IWSConn) {
				called = true
			},
		}

		app := fibertools.NewApp(
			PluginFiveInARow(roomManager),
		)

		port := fmt.Sprintf("%d", freePort())
		go app.Listen(":" + port)
		defer app.Shutdown()

		client := &http.Client{}
		req := gotils.ResultOrPanic(http.NewRequest(
			"GET",
			"http://127.0.0.1:"+port+path,
			&bytes.Buffer{},
		))
		req.Header.Add("Connection", "Upgrade")
		req.Header.Add("Upgrade", "Websocket")
		req.Header.Add("Sec-Websocket-Version", "13")
		req.Header.Add("Sec-WebSocket-Key", "fake-key")
		var resp *http.Response
		var err error

		for {
			resp, err = client.Do(req)
			if err == nil {
				break
			}
			time.Sleep(1 * time.Millisecond)
		}

		time.Sleep(1 * time.Millisecond)
		assert.Equal(t, http.StatusSwitchingProtocols, resp.StatusCode)
		assert.True(t, called)
	})
}

func freePort() int {
	addr := gotils.ResultOrPanic(net.ResolveTCPAddr("tcp", "localhost:0"))
	listener := gotils.ResultOrPanic(net.ListenTCP("tcp", addr))
	defer listener.Close()
	return listener.Addr().(*net.TCPAddr).Port
}
