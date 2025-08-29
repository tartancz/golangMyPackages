package discord

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"strings"
	"time"
)

type TCPTransport struct {
	parser      MessageParser
	router      MessageRouter
	onRaw       func(string, io.Writer)
	delimiter   string
	conn        net.Conn
	messageChan chan []byte
	config      Config
}

// Config for TCPTransport.
type Config struct {
	ProgramName string
	Host        string
	Port        string
	DialTimeout time.Duration
	RetryDelay  time.Duration
}

func NewTCPTransport(cfg Config, parser MessageParser, router MessageRouter) *TCPTransport {
	return &TCPTransport{
		parser:      parser,
		router:      router,
		config:      cfg,
		messageChan: make(chan []byte),
	}
}

func (t *TCPTransport) SetOnRawMessage(f func(string, io.Writer)) {
	t.onRaw = f
}

func (t *TCPTransport) Start() error {
	addr := net.JoinHostPort(t.config.Host, t.config.Port)
	var err error

	timeout := time.NewTimer(t.config.DialTimeout)
	defer timeout.Stop()

	for {
		select {
		case <-time.After(t.config.RetryDelay):
			t.conn, err = net.Dial("tcp", addr)
			if err == nil {
				goto connected
			}
			fmt.Println("Retrying connection:", err)
		case <-timeout.C:
			return fmt.Errorf("timeout connecting to %s", addr)
		}
	}

connected:
	fmt.Fprintf(t.conn, "SET_NAME:%s\n", t.config.ProgramName)
	fmt.Fprintf(t.conn, "GET_DELIMITER:\n")

	go t.readLoop()
	return t.writeLoop()
}

func (t *TCPTransport) readLoop() {
	reader := bufio.NewReader(t.conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		if strings.HasPrefix(line, "DELIMITER:") {
			t.delimiter = strings.TrimSpace(strings.TrimPrefix(line, "DELIMITER:"))
			continue
		}

		if t.onRaw != nil {
			t.onRaw(line, t)
		}

		msg, err := t.parser.Parse(line)
		if err != nil {
			fmt.Fprintln(t, "Parse error:", err)
			continue
		}

		if t.parser.IsHelpRequest(msg) {
			t.router.Help(t)
		} else {
			t.router.Route(msg, t)
		}
	}
}

func (t *TCPTransport) writeLoop() error {
	for msg := range t.messageChan {
		if _, err := t.conn.Write(msg); err != nil {
			return err
		}
	}
	return nil
}

func (t *TCPTransport) Write(p []byte) (int, error) {
	if t.conn == nil {
		return 0, fmt.Errorf("no connection")
	}
	msg := append(bytes.TrimSpace(p), []byte("\n"+t.delimiter+"\n")...)
	t.messageChan <- msg
	return len(p), nil
}

func (t *TCPTransport) Close() {
	if t.messageChan != nil {
		close(t.messageChan)
	}
	if t.conn != nil {
		t.conn.Close()
	}
}
