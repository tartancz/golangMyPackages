package discord

import "io"

// Message represents a parsed incoming command.
type Message struct {
	Raw         string
	CommandName string
	Args        []string
}

// Handler processes a command.
type Handler interface {
	Handle(msg Message, out io.Writer)
	GetHelp() string
}

// MessageParser turns raw input into structured commands.
type MessageParser interface {
	Parse(raw string) (Message, error)
	IsHelpRequest(msg Message) bool
}

// MessageRouter dispatches a message to the appropriate handler.
type MessageRouter interface {
	Route(msg Message, out io.Writer)
	Help(out io.Writer)
	Register(command string, handler Handler)
}

// Server coordinates input/output and routing.
type Server interface {
	SetOnRawMessage(func(string, io.Writer))
	Start() error
	Close()
}
