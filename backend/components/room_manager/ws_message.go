package room_manager

import (
	"encoding/json"

	"github.com/DAtek/golidator"
)

type MessageType string

type messageStruct struct {
	Type MessageType `json:"type"`
	Data string      `json:"data"`
}

func (msg *messageStruct) setData(data any, msgType MessageType) {
	msg.Type = msgType
	d, _ := json.Marshal(data)
	msg.Data = string(d)
}

var okMessage = messageStruct{
	Type: MessageTypeOK,
}

const (
	MessageTypeOK          = MessageType("OK")
	MessageTypeBadMessage  = MessageType("BAD_MESSAGE")
	MessageTypeJoin        = MessageType("JOIN")
	MessageTypePickSide    = MessageType("PICK_SIDE")
	MessageTypeMove        = MessageType("MOVE")
	MessageTypeUpdateGame  = MessageType("UPDATE_GAME")
	MessageTypeSendMessage = MessageType("SEND_MESSAGE")
)

func (obj *messageStruct) GetValidators(params ...interface{}) golidator.ValidatorCollection {
	return golidator.ValidatorCollection{
		{Field: "type", Function: func() *golidator.ValueError {
			for _, msgType := range getInputMessageTypes() {
				if msgType == obj.Type {
					return nil
				}
			}
			return &golidator.ValueError{ErrorType: "INVALID_MESSAGE_TYPE"}
		}},
	}
}

var inputMessageTypes = []MessageType{}

func getInputMessageTypes() []MessageType {
	if len(inputMessageTypes) > 0 {
		return inputMessageTypes
	}

	for key := range actions {
		inputMessageTypes = append(inputMessageTypes, key)
	}

	return inputMessageTypes
}

func parseData[T interface{}](obj *T, msg *messageStruct) (*T, *string) {
	if msg.Data == "" {
		return nil, &invalidDataError
	}

	if err := json.Unmarshal([]byte(msg.Data), obj); err != nil {
		return nil, &invalidDataError
	}

	return obj, nil
}

var invalidDataError = "INVALID_DATA"
