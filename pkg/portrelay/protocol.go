package portrelay

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

type MessageProtocol interface {
	Encode(message Message) []byte
	Decode(reader io.Reader) (*Message, error)
}

type Message struct {
	Command   string
	Arguments []string
}

// FORMAT
// BasicMessageProtocol: "*<number of arguments>\n$<number of bytes of argument 1>\n<argument data>\n..."
type BinaryMessageProtocol struct{}

func NewBinaryMessageProtocol() *BinaryMessageProtocol {
	return &BinaryMessageProtocol{}
}

func (p *BinaryMessageProtocol) Encode(message Message) []byte {

	var builder strings.Builder

	builder.WriteString(fmt.Sprintf("*%d\n", len(message.Arguments)+1))

	args := append([]string{message.Command}, message.Arguments...)
	for _, arg := range args {
		builder.WriteString(fmt.Sprintf("$%d\n", len(arg)))
		builder.WriteString(arg)
		builder.WriteString("\n")
	}

	return []byte(builder.String())
}

func (p *BinaryMessageProtocol) DecodeString(s string) (*Message, error) {
	return p.Decode(strings.NewReader(s))
}

func (p *BinaryMessageProtocol) DecodeBytes(b []byte) (*Message, error) {
	return p.Decode(bytes.NewReader(b))
}

func (p *BinaryMessageProtocol) Decode(reader io.Reader) (*Message, error) {
	var msg Message

	buf := bufio.NewReader(reader)

	// Read the first line: "*<number of args>\n"
	line, err := buf.ReadString('\n')
	if err != nil {
		return nil, &DecodeError{
			Stage:   "read argument count line",
			Index:   -1,
			Details: "could not read '*<n>' line",
			Err:     err,
		}
	}

	var lenArgs int
	if _, err := fmt.Sscanf(line, "*%d\n", &lenArgs); err != nil {
		return nil, &DecodeError{
			Stage:   "parse argument count",
			Index:   -1,
			Details: fmt.Sprintf("invalid line: %q", strings.TrimSpace(line)),
			Err:     err,
		}
	}

	args := make([]string, lenArgs)

	for i := 0; i < lenArgs; i++ {
		// Read: "$<length>\n"
		line, err := buf.ReadString('\n')
		if err != nil {
			return nil, &DecodeError{
				Stage:   "read argument length line",
				Index:   i,
				Details: "could not read '$<length>' line",
				Err:     err,
			}
		}

		var argLen int
		if _, err := fmt.Sscanf(line, "$%d\n", &argLen); err != nil {
			return nil, &DecodeError{
				Stage:   "parse argument length",
				Index:   i,
				Details: fmt.Sprintf("invalid line: %q", strings.TrimSpace(line)),
				Err:     err,
			}
		}

		// Read actual argument data
		argData := make([]byte, argLen)
		n, err := io.ReadFull(buf, argData)
		if err != nil || n != argLen {
			return nil, &DecodeError{
				Stage:   "read argument data",
				Index:   i,
				Details: fmt.Sprintf("expected %d bytes, got %d", argLen, n),
				Err:     err,
			}
		}
		args[i] = string(argData)

		// Expect newline after data
		newline, err := buf.ReadString('\n')
		if err != nil {
			return nil, &DecodeError{
				Stage:   "read newline after data",
				Index:   i,
				Details: "could not read expected newline after argument",
				Err:     err,
			}
		}
		if newline != "\n" {
			return nil, &DecodeError{
				Stage:   "validate newline after data",
				Index:   i,
				Details: fmt.Sprintf("expected newline, got %q", strings.TrimRight(newline, "\r\n")),
				Err:     fmt.Errorf("invalid format"),
			}
		}
	}

	msg.Command = args[0]
	msg.Arguments = args[1:]

	return &msg, nil
}
