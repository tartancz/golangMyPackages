package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	l, e := net.DialTimeout("tcp", net.JoinHostPort("localhost", "8080"), 60*time.Second)
	fmt.Print(l, e)
	// p := portrelay.NewBinaryMessageProtocol()
	// r := p.Encode(portrelay.Message{
	// 	Command:   "TestCommand",
	// 	Arguments: []string{"-t", "TestArgument"},
	// })

	// fmt.Printf("%s", strconv.Quote(string(r)))
}
