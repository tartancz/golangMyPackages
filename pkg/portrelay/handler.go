package portrelay

import "io"

type Handler interface {
	Handle(msg Message, out io.Writer)
	GetHelp() string
}

type FuncHandler struct {
	Func func(msg Message, out io.Writer)
	Help string
}

func (f FuncHandler) Handle(msg Message, out io.Writer) {
	f.Func(msg, out)
}

func (f FuncHandler) GetHelp() string {
	if f.Help == "" {
		return "No help available."
	}
	return f.Help
}
