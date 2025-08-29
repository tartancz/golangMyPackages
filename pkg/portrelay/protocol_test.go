package portrelay

import (
	"bytes"
	"reflect"
	"strconv"
	"testing"
)

func TestDecode(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected *Message
		wantErr  bool
	}{
		{name: "Valid Command without arguments", input: []byte("*1\n$11\nTestCommand\n"), expected: &Message{Command: "TestCommand", Arguments: []string{}}, wantErr: false},
		{name: "Valid Command with arguments", input: []byte("*3\n$11\nTestCommand\n$2\n-t\n$12\nTestArgument\n"), expected: &Message{Command: "TestCommand", Arguments: []string{"-t", "TestArgument"}}, wantErr: false},
		{name: "Invalid Lenght of arguments format", input: []byte("*5\n$11\nTestCommand\n$2\n-t\n"), expected: nil, wantErr: true},
		{name: "Invalid Argument Data format", input: []byte("*3\n$11\nTestCommand\n$2\n-t\n$12\nTestArgument"), expected: nil, wantErr: true},
		{name: "Empty Input", input: []byte(""), expected: nil, wantErr: true},
		{name: "Invalid Input", input: []byte("Invalid Input"), expected: nil, wantErr: true},
		{name: "Empty Command", input: []byte("*1\n$0\n\n"), expected: &Message{Command: "", Arguments: []string{}}, wantErr: false},
		{name: "Command with empty argument", input: []byte("*2\n$11\nTestCommand\n$0\n\n"), expected: &Message{Command: "TestCommand", Arguments: []string{""}}, wantErr: false},
		{name: "Command with spaces", input: []byte("*3\n$12\nTest Command\n$5\narg 1\n$5\narg 2\n"), expected: &Message{Command: "Test Command", Arguments: []string{"arg 1", "arg 2"}}, wantErr: false},
	}
	
	p := NewBinaryMessageProtocol()

	for _, tt := range tests {

		t.Run("READER: " + tt.name, func(t *testing.T) {
			
			got, err := p.Decode(bytes.NewReader(tt.input))
			if (err != nil) != tt.wantErr {
				t.Errorf("Decode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("Decode() got = %v, want %v", got, tt.expected)
			}
		})

		t.Run("STRING: " + tt.name, func(t *testing.T) {
			got, err := p.DecodeString(string(tt.input))
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("DecodeString() got = %v, want %v", got, tt.expected)
			}
		})

		t.Run("BYTES: " + tt.name, func(t *testing.T) {
			got, err := p.DecodeBytes(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("DecodeBytes() got = %v, want %v", got, tt.expected)
			}
		})
	}


}



func TestEncode(t *testing.T) {
	tests := []struct {
		name     string
		message  *Message
		expected []byte
	}{
		{name: "Command without arguments", message: &Message{Command: "TestCommand", Arguments: []string{}}, expected: []byte("*1\n$11\nTestCommand\n")},
		{name: "Command with arguments", message: &Message{Command: "TestCommand", Arguments: []string{"-t", "TestArgument"}}, expected: []byte("*3\n$11\nTestCommand\n$2\n-t\n$12\nTestArgument\n")},
		{name: "Command With spaces", message: &Message{Command: "Test Command", Arguments: []string{"arg 1", "arg 2"}}, expected: []byte("*3\n$12\nTest Command\n$5\narg 1\n$5\narg 2\n")},
		{name: "Empty string", message: &Message{}, expected: []byte("*1\n$0\n\n")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewBinaryMessageProtocol()
			got := p.Encode(*tt.message)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("Encode() got = %v (%v), want %v (%v)", got, strconv.Quote(string(got)), tt.expected, strconv.Quote(string(tt.expected)))
			}
		})

	}

}
