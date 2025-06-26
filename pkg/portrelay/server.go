package portrelay

import (
	"errors"
	"io"
	"net"
	"os"
)

type MessageArgs struct {
	Raw         string
	Args        []string
	CommandName string
}

type MessageHandler interface {
	HandleMessage(args MessageArgs, writer io.Writer)
	GetHelpMessage() string
}

type Server struct {
	messageChan chan string
	onMessage   func(string, io.Writer)
	handlers    map[string]MessageHandler
	conn        net.Conn
	protocol    MessageProtocol
}

func NewServer(protocol MessageProtocol) *Server {
	return &Server{
		messageChan: make(chan string),
		onMessage:   nil,
		handlers:    make(map[string]MessageHandler),
		conn:        nil,
		protocol:    protocol,
	}
}

func (s *Server) Start(host, port string) error {
	if s.conn != nil {
		return errors.New("server already started")
	}

	if host == "" {
		if host, exists := os.LookupEnv("BOT_SERVER_HOST"); exists {
			host = host
		} else {
			return errors.New("no host specified")
		}
	}

	if port == "" {
		if port, exists := os.LookupEnv("BOT_SERVER_PORT"); exists {
			port = port
		} else {
			return errors.New("no port specified")
		}
	}

	l, err := net.Dial("tcp", net.JoinHostPort(host, port))
	if err != nil {
		return err
	}
	s.conn = l
}
