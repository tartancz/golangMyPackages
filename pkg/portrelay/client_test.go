package portrelay

import (
	"io"
	"net"
	"strings"
	"testing"
)

func pingHandler(msg Message, out io.Writer) {
	if strings.ToLower(msg.Command) != "ping" {
		return
	}
	out.Write([]byte("pong"))
}

func TestClientStart_MockDial(t *testing.T) {
	client := NewClient(NewBinaryMessageProtocol())

	writer := make(chan []byte)
	defer close(writer)

	client.dial = func(network, address string) (net.Conn, error) {
		c1, c2 := net.Pipe() 
		go func() {
			defer c2.Close()
			
			for msg := range writer {
				_, err := c2.Write(msg)
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}
		}()
		return c1, nil
	}

	err := client.Start("fakehost", "1234")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	defer client.Close()
	
	tests := []struct {
		name     string
		input    []byte
		expected *Message
	}{
		{}
	}


}

// func TestClientHandler(t *testing.T) {
// 	buffer := bytes.NewBuffer(nil)
// 	protocol := NewBinaryMessageProtocol()
// 	client := NewClient(protocol)
// 	client
// }
