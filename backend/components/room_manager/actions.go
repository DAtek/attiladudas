package room_manager

var actions = actionCollection{
	MessageTypeJoin:        joinRoom,
	MessageTypePickSide:    pickSide,
	MessageTypeSendMessage: sendMessage,
	MessageTypeMove:        move,
}
