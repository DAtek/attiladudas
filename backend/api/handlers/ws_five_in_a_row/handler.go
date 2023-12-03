package fiar

import (
	"api/components/room_manager"
	"fibertools"
	"fmt"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

const path = "/ws/five-in-a-row/"

func PluginFiveInARow(roomManager room_manager.IRoomManager) fibertools.Plugin {
	cfg := websocket.Config{
		RecoverHandler: func(conn *websocket.Conn) {
			if err := recover(); err != nil {
				conn.WriteJSON(fiber.Map{"error": fmt.Sprintf("%v", err)})
			}
		},
	}

	return func(app *fiber.App) {
		app.Get(
			path,
			websocket.New(
				func(c *websocket.Conn) {
					roomManager.HandleConnection(c)
				},
				cfg,
			),
		)
	}
}
