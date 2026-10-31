// Copyright (c) 2026 Parable, LLC
package protocol

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestFrames(t *testing.T) {
	for _, input := range []string{`{}`, "{} {}\n", strings.Repeat("x", MaxMessage) + "\n", "{\"unknown\":1}\n"} {
		var r Request
		if Read(strings.NewReader(input), &r) == nil {
			t.Fatalf("accepted invalid frame %q", input[:min(30, len(input))])
		}
	}
	var b bytes.Buffer
	r := Request{Version: 1, ID: "test", Method: "service.status"}
	if err := Write(&b, r); err != nil {
		t.Fatal(err)
	}
	var got Request
	if err := Read(&b, &got); err != nil || got != r {
		t.Fatalf("round trip: %+v %v", got, err)
	}
}
func TestGolden(t *testing.T) {
	b, err := os.ReadFile("../../tests/fixtures/protocol/status.json")
	if err != nil {
		t.Fatal(err)
	}
	var r Response
	if err = Read(bytes.NewReader(b), &r); err != nil {
		t.Fatal(err)
	}
	if r.Result.Run.State != "succeeded" || *r.Result.Run.ExitCode != 0 {
		t.Fatal(r)
	}
}
func TestCompatibility(t *testing.T) {
	if e := Validate(Request{Version: 2}); e == nil || e.Code != "incompatible_version" {
		t.Fatal(e)
	}
}
