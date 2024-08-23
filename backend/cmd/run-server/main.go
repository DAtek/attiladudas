package main

import (
	"api"
	"api/components/room_manager"
	fiar "api/handlers/ws_five_in_a_row"
	"flag"
	"net"

	"github.com/DAtek/gotils"
	"github.com/gofiber/fiber/v2"
)

func main() {
	port := flag.String("port", "8000", "-port 8000")
	host := flag.String("host", "127.0.0.1", "-host 127.0.0.1")
	socket := flag.String("socket", "", "-socket /var/run/attiladudas/server.sock")
	flag.Parse()
	app := createApp()
	listener := getListener(*host, *port, *socket)
	app.Listener(listener)
}

func createApp() *fiber.App {
	roomManager := room_manager.NewRoomManager()

	return api.AppWithMiddlewares(
		fiar.PluginFiveInARow(roomManager),
	)
}

func getListener(host, port, socket string) net.Listener {
	if socket == "" {
		addr := gotils.ResultOrPanic(
			net.ResolveTCPAddr("tcp4", host+":"+port),
		)
		return gotils.ResultOrPanic(net.ListenTCP("tcp4", addr))
	}

	addr := gotils.ResultOrPanic(
		net.ResolveUnixAddr("unix", socket),
	)
	return gotils.ResultOrPanic(net.ListenUnix("unix", addr))
}
