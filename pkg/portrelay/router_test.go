package portrelay

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestCommandRouter_Route_KnownCommand(t *testing.T) {
	var output bytes.Buffer

	router := NewRouter()
	router.Register("ping", FuncHandler{
		Func: func(msg Message, out io.Writer) {
			out.Write([]byte("pong\n"))
		},
		Help: "Responds with pong",
	})

	msg := Message{Command: "ping"}
	router.Route(msg, &output)

	expected := "pong\n"
	if output.String() != expected {
		t.Errorf("Expected %q but got %q", expected, output.String())
	}
}

func TestCommandRouter_Route_UnknownCommand(t *testing.T) {
	var output bytes.Buffer

	router := NewRouter()
	router.Register("ping", FuncHandler{
		Func: func(msg Message, out io.Writer) {
			out.Write([]byte("pong\n"))
		},
		Help: "Responds with pong",
	})

	msg := Message{Command: "unknown"}
	router.Route(msg, &output)

	outStr := output.String()
	if !strings.Contains(outStr, "Unknown command: unknown") {
		t.Errorf("Expected unknown command message, got: %q", outStr)
	}
	if !strings.Contains(outStr, "ping: Responds with pong") {
		t.Errorf("Expected help output to include registered commands, got: %q", outStr)
	}
}

func TestCommandRouter_Help_WithCommands(t *testing.T) {
	var output bytes.Buffer

	router := NewRouter()
	router.Register("hello", FuncHandler{
		Func: func(msg Message, out io.Writer) {
			out.Write([]byte("hi there\n"))
		},
		Help: "Says hello",
	})

	router.Help(&output)

	outStr := output.String()
	if !strings.Contains(outStr, "hello: Says hello") {
		t.Errorf("Expected help to include 'hello' command, got: %q", outStr)
	}
}

func TestCommandRouter_Help_NoCommands(t *testing.T) {
	var output bytes.Buffer

	router := NewRouter()
	router.Help(&output)

	expected := "No commands available.\n"
	if output.String() != expected {
		t.Errorf("Expected %q but got %q", expected, output.String())
	}
}

func TestCommandRouter_CommandCaseInsensitive(t *testing.T) {
	var output bytes.Buffer

	router := NewRouter()
	router.Register("HeLp", FuncHandler{
		Func: func(msg Message, out io.Writer) {
			out.Write([]byte("show help\n"))
		},
		Help: "Shows help",
	})

	msg := Message{Command: "help"}
	router.Route(msg, &output)

	expected := "show help\n"
	if output.String() != expected {
		t.Errorf("Expected %q but got %q", expected, output.String())
	}
}
