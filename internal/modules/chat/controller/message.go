package controller

// MessageController handles incoming HTTP/WebSocket requests for chat messages.
type MessageController struct {
}

func NewMessageController() *MessageController {
	return &MessageController{}
}
