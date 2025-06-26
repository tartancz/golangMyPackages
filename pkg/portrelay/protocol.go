package portrelay

type MessageProtocol interface {
	Encode(message string) []byte
	Decode(message []byte) string
}


type DiscordMessageProtocol struct {
	delimiter string

}

func NewDiscordMessageProtocol(delimiter string) *DiscordMessageProtocol {
	return &DiscordMessageProtocol{delimiter: delimiter}
}

func (p *DiscordMessageProtocol) Encode(message string) []byte {
	return []byte(message + p.delimiter)
}

func (p *DiscordMessageProtocol) Decode(message []byte) string {
	return string(message[:len(message)-len(p.delimiter)])
}