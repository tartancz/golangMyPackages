package discord

import (
	"fmt"
	"io"
	"strings"
)

type CommandRouter struct {
	handlers map[string]Handler
}

func NewRouter() *CommandRouter {
	return &CommandRouter{
		handlers: make(map[string]Handler),
	}
}

func (r *CommandRouter) Register(command string, handler Handler) {
	r.handlers[strings.ToLower(command)] = handler
}

func (r *CommandRouter) Route(msg Message, out io.Writer) {
	handler, ok := r.handlers[msg.CommandName]
	if !ok {
		fmt.Fprintf(out, "Unknown command: %s\n\n", msg.CommandName)
		r.Help(out)
		return
	}
	handler.Handle(msg, out)
}

func (r *CommandRouter) Help(out io.Writer) {
	if len(r.handlers) == 0 {
		fmt.Fprintln(out, "No commands available.")
		return
	}
	for name, h := range r.handlers {
		fmt.Fprintf(out, "%s: %s\n-----------------\n", name, h.GetHelp())
	}
}
