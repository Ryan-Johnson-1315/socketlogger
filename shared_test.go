package socketlogger

import "testing"

func TestSocketMessageType(t *testing.T) {
	var msg SocketMessage = &LogMessage{}
	if msg.Type() != Log {
		t.Errorf("msg should be Log type, actual: %T", msg)
	}
}
