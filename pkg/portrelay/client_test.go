package portrelay

import (
	"bytes"
	"testing"
)

func TestClientHandler(t *testing.T) {
	buffer := bytes.NewBuffer(nil)
	protocol := NewBinaryMessageProtocol()
	client := NewClient(protocol)
	client
}
