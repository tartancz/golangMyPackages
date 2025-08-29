package main

/*
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"

	"github.com/tartancz/golangMyPackages/pkg/portrelay"
)


//export EncodeMessage
func EncodeMessage(command *C.char, argc C.int, argv **C.char) *C.char {
	args := make([]string, int(argc))
	argSlice := (*[1 << 30]*C.char)(unsafe.Pointer(argv))[:argc:argc]
	for i, a := range argSlice {
		args[i] = C.GoString(a)
	}

	msg := portrelay.Message{
		Command:   C.GoString(command),
		Arguments: args,
	}

	proto := portrelay.NewBinaryMessageProtocol()
	encoded := proto.Encode(msg)
	return C.CString(string(encoded))
}

//export DecodeMessage
func DecodeMessage(data *C.char) *C.char {
	s := C.GoString(data)
	proto := portrelay.NewBinaryMessageProtocol()
	msg, err := proto.DecodeString(s)
	if err != nil {
		return C.CString("ERROR:" + err.Error())
	}

	// For simplicity return JSON
	json := `{"command":"` + msg.Command + `","args":[`
	for i, a := range msg.Arguments {
		if i > 0 {
			json += ","
		}
		json += `"` + a + `"`
	}
	json += "]}"
	return C.CString(json)
}
