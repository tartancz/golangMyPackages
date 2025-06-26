package discord

import (
	"errors"
	"strings"
)

type DefaultParser struct{}

func (p *DefaultParser) Parse(raw string) (Message, error) {
	raw = strings.TrimSpace(raw)
	tokens := strings.Fields(raw)

	if len(tokens) < 3 {
		return Message{}, errors.New("invalid message format")
	}

	return Message{
		Raw:         raw,
		CommandName: strings.ToLower(tokens[2]),
		Args:        tokens[3:],
	}, nil
}

func (p *DefaultParser) IsHelpRequest(msg Message) bool {
	return msg.CommandName == "help"
}
