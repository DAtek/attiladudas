package room_manager

import (
	"encoding/json"
	"log"

	"github.com/DAtek/golidator"
	"github.com/gorilla/websocket"
)

func (manager *roomManager) HandleConnection(conn IWSConn) {
	defer func() {
		manager.cleanup_(conn)
		log.Println("Connection closed")
	}()

	for {
		messageType, data, err := conn.ReadMessage()
		if err != nil {
			return
		}

		if messageType == websocket.CloseMessage {
			return
		}

		gameMessage := &messageStruct{}
		json.Unmarshal(data, gameMessage)
		validationErr := golidator.Validate(gameMessage)

		if handleBadMessage(conn, validationErr, messageType) {
			continue
		}

		action := manager.actions[MessageType(gameMessage.Type)]
		msg := action(manager, conn, gameMessage)
		log.Printf("Game action: %s, result: %s\n", gameMessage.Type, msg.Type)
		msgData, _ := json.Marshal(msg)
		conn.WriteMessage(messageType, msgData)
	}
}
