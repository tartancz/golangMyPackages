package portrelay

import "fmt"

type DecodeError struct {
	Stage   string // e.g. "read argument length", "parse argument data"
	Index   int    // argument index (or -1 if not applicable)
	Details string // optional extra context
	Err     error  // wrapped error
}

func (e *DecodeError) Error() string {
	return fmt.Sprintf("decode error at stage %q (arg %d): %s: %v", e.Stage, e.Index, e.Details, e.Err)
}

func (e *DecodeError) Unwrap() error {
	return e.Err
}

// throws whenever the server cannot connect to the specified host and port
// this is used to distinguish between connection errors and other types of errors
type ConnError struct {
	Err     error
}

func (e *ConnError) Error() string {
	return fmt.Sprintf("failed to connect to server: %v", e.Err)
}

func (e *ConnError) Unwrap() error {
	return e.Err
}
