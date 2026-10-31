// Copyright (c) 2026 Parable, LLC
package protocol

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"io"
)

const Version = 1
const MaxMessage = 64 * 1024

type Request struct {
	Version int    `json:"version"`
	ID      string `json:"id"`
	Method  string `json:"method"`
	Fixture string `json:"fixture,omitempty"`
}
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
type Run struct {
	State    string `json:"state"`
	Output   string `json:"output"`
	ExitCode *int   `json:"exit_code,omitempty"`
}
type Result struct {
	PID int  `json:"pid"`
	Run *Run `json:"run,omitempty"`
}
type Response struct {
	Version int     `json:"version"`
	ID      string  `json:"id"`
	Result  *Result `json:"result,omitempty"`
	Error   *Error  `json:"error,omitempty"`
}

func Failure(id, code, message string) Response {
	return Response{Version: Version, ID: id, Error: &Error{code, message}}
}

// Read accepts exactly one bounded newline-delimited JSON frame. A connection
// carries one request and one response; trailing frames are never executed.
func Read(r io.Reader, v any) error {
	b, err := bufio.NewReader(io.LimitReader(r, MaxMessage+1)).ReadBytes('\n')
	if len(b) > MaxMessage {
		return errors.New("message exceeds 65536 bytes")
	}
	if err != nil {
		return errors.New("message requires newline terminator")
	}
	d := json.NewDecoder(bytes.NewReader(b))
	d.DisallowUnknownFields()
	if err = d.Decode(v); err != nil {
		return err
	}
	var extra any
	if d.Decode(&extra) != io.EOF {
		return errors.New("trailing JSON")
	}
	return nil
}
func Write(w io.Writer, v any) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	if len(b)+1 > MaxMessage {
		return errors.New("message exceeds limit")
	}
	_, err = w.Write(append(b, '\n'))
	return err
}
func Validate(r Request) *Error {
	if r.Version != Version {
		return &Error{"incompatible_version", "Expected protocol version 1"}
	}
	if len(r.ID) == 0 || len(r.ID) > 128 {
		return &Error{"invalid_request", "Request id must contain 1–128 bytes"}
	}
	if r.Method != "service.status" && r.Method != "fixture.start" && r.Method != "fixture.cancel" {
		return &Error{"invalid_request", "Unknown method"}
	}
	if r.Method == "fixture.start" && r.Fixture != "complete" && r.Fixture != "cancel" {
		return &Error{"invalid_request", "Unknown fixture"}
	}
	if r.Method != "fixture.start" && r.Fixture != "" {
		return &Error{"invalid_request", "Unexpected fixture"}
	}
	return nil
}
