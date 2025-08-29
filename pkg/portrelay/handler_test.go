package portrelay

import (
	"bytes"
	"io"
	"testing"
)

func TestFuncHandler_Handle(t *testing.T) {
	var buf bytes.Buffer
	handler := FuncHandler{
		Func: func(msg Message, out io.Writer) {
			out.Write([]byte("handled: " + msg.Command))
		},
		Help: "test handler",
	}

	msg := Message{Command: "test"}
	handler.Handle(msg, &buf)

	expected := "handled: test"
	if buf.String() != expected {
		t.Errorf("Expected output %q but got %q", expected, buf.String())
	}
}

func TestFuncHandler_GetHelp_WithHelp(t *testing.T) {
	handler := FuncHandler{
		Func: func(msg Message, out io.Writer) {},
		Help: "this is a handler",
	}

	help := handler.GetHelp()
	expected := "this is a handler"
	if help != expected {
		t.Errorf("Expected help %q but got %q", expected, help)
	}
}

func TestFuncHandler_GetHelp_NoHelp(t *testing.T) {
	handler := FuncHandler{
		Func: func(msg Message, out io.Writer) {},
		Help: "",
	}

	help := handler.GetHelp()
	expected := "No help available."
	if help != expected {
		t.Errorf("Expected help %q but got %q", expected, help)
	}
}
