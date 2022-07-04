package event

import (
	"encoding/json"
	"fmt"
)

const MessageEvent = "message_event"

func NewMessage(sender, message string) string {
	b, _ := json.Marshal(struct {
		Type    string `json:"type"`
		Content string `json:"content"`
	}{
		Type:    MessageEvent,
		Content: fmt.Sprintf("%s: %s", sender, message),
	})

	return string(b)
}
