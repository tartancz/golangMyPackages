package portrelay

import (
	"errors"
	"io"
	"net"
	"os"
)

type Client struct {
	//
	messageChan  chan Message
	onAnyMessage func(string, io.Writer)
	onUnhandled func(Message, io.Writer)
	handlers     map[string]Handler
	conn         net.Conn
	protocol     MessageProtocol
}

func NewServer(protocol MessageProtocol) *Client {
	return &Client{
		//		messageChan: make(chan string),
		onAnyMessage: nil,
		handlers:     make(map[string]Handler),
		conn:         nil,
		protocol:     protocol,
	}
}

func (s *Client) Start(host, port string) error {
	if s.conn != nil {
		return errors.New("server already started")
	}

	if host == "" {
		var exists bool
		if host, exists = os.LookupEnv("BOT_SERVER_HOST"); !exists {
			return errors.New("no host specified")
		}
	}

	if port == "" {
		var exists bool
		if port, exists = os.LookupEnv("BOT_SERVER_PORT"); !exists {
			return errors.New("no port specified")
		}
	}

	c, err := net.Dial("tcp", net.JoinHostPort(host, port))
	if err != nil {
		return &ConnError{Err: err} //errors.New("failed to connect to server: " + err.Error())
	}
	s.conn = c
	s.messageChan = make(chan Message)

	go func() {
		for msg := range s.messageChan {
			encoded := s.protocol.Encode(msg)
			_, err := c.Write(encoded)
			if err != nil {
				return
			}
	}}()
	
	go func() {
		for {
			message, err := s.protocol.Decode(c)
			if err != nil {
				//TODO: better error handling
				return
			}

			if handler, exists := s.handlers[message.Command]; exists {
				go handler.Handle(*message, c)
			} else if s.onUnhandled != nil {
				go s.onUnhandled(*message, c)
			}

			if s.onAnyMessage != nil {
				go s.onAnyMessage(message.Command, c)
			}
		}
	}()

	s.conn = c
	return nil
}

func (s *Client) StartWithRetry(host, port string, retries int) error {
	for i := 0; i < retries; i++ {
		var connErr *ConnError
		if err := s.Start(host, port); errors.Is(err, connErr) {
			continue
		} else {
			return err
		}

	}

	return errors.New("failed to start server after retries")
}


func (c *Client) SendMessage(msg Message) {
	c.messageChan <- msg
}