package portrelay

import (
	"errors"
	"io"
	"net"
	"os"
	"strings"
)

type Dialer func(network, address string) (net.Conn, error)

type Client struct {
	//
	messageChan  chan Message
	OnAnyMessage func(string, io.Writer)
	OnUnhandled  func(Message, io.Writer)
	Handlers     map[string]Handler
	conn         net.Conn
	protocol     MessageProtocol
	dial         Dialer
}

func NewClient(protocol MessageProtocol) *Client {
	return &Client{
		Handlers: make(map[string]Handler),
		protocol: protocol,
		dial:     net.Dial,
	}
}

func (s *Client) Start(host, port string) error {
	if s.conn != nil {
		return errors.New("client already started")
	}

	if host == "" {
		var exists bool
		if host, exists = os.LookupEnv("BOT_CLIENT_HOST"); !exists {
			return errors.New("no host specified")
		}
	}

	if port == "" {
		var exists bool
		if port, exists = os.LookupEnv("BOT_CLIENT_PORT"); !exists {
			return errors.New("no port specified")
		}
	}

	c, err := s.dial("tcp", net.JoinHostPort(host, port))
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
		}
	}()

	go func() {
		for {
			message, err := s.protocol.Decode(c)
			if err != nil {
				//TODO: better error handling
				return
			}

			if handler, exists := s.Handlers[message.Command]; exists {
				go handler.Handle(*message, c)
			} else if s.OnUnhandled != nil {
				go s.OnUnhandled(*message, c)
			}

			if s.OnAnyMessage != nil {
				go s.OnAnyMessage(message.Command, c)
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

func (c *Client) RegisterHandler(command string, handler Handler) {
	c.Handlers[strings.ToLower(command)] = handler
}

func (c *Client) Close() error {
	return nil
}
