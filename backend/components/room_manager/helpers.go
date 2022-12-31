package room_manager

import (
	"attiladudas/backend/helpers"
	"attiladudas/backend/ws"
	"encoding/json"

	"github.com/DAtek/golidator"
)

func handleBadMessage(conn ws.IConn, validationError *golidator.ValidationError, messageType int) bool {
	if validationError == nil {
		return false
	}

	msg := createBadMessageFromValidationError(validationError)
	responseData, _ := json.Marshal(msg)
	conn.WriteMessage(messageType, responseData)
	return true
}

func createBadMessageFromValidationError(validationError *golidator.ValidationError) messageStruct {
	jsonError := helpers.JsonErrorFromValidationError(validationError)
	errData, _ := json.Marshal(jsonError)
	errStr := string(errData)
	return messageStruct{Type: MessageTypeBadMessage, Data: errStr}
}
