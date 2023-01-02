package main

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"
	"time"
)

func TestGetNameContext(t *testing.T) {
	cases := []struct {
		name         string
		timeout      time.Duration
		input        string
		output       string
		expectedName string
		timedOut     bool
		err          error
	}{
		{
			name:         "GetNameContextWithEmptyInput",
			timeout:      100 * time.Millisecond,
			input:        "",
			output:       "Your name please? Press the key Enter when done\n",
			expectedName: "",
			timedOut:     false,
			err:          errors.New("you entered an empty name"),
		},
		{
			name:         "GetNameContextTimedOut",
			timeout:      100 * time.Millisecond,
			input:        "",
			output:       "Your name please? Press the key Enter when done\n",
			expectedName: "Default Name",
			timedOut:     true,
			err:          context.DeadlineExceeded,
		},
		{
			name:         "GetNameContextWithoutTimeout",
			timeout:      100 * time.Millisecond,
			input:        "John Doe",
			output:       "Your name please? Press the key Enter when done\n",
			expectedName: "John Doe",
			timedOut:     false,
			err:          nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
			defer cancel()
			var r io.Reader
			r = strings.NewReader(c.input)
			if c.timedOut {
				r = &TimedReader{duration: 200 * time.Millisecond}
			}
			w := bytes.Buffer{}
			name, err := getNameContext(ctx, r, &w)
			if c.err == nil && err != nil {
				t.Fatalf("expected nil error, got %q", err)
			}
			if c.err != nil && c.err.Error() != err.Error() {
				t.Fatalf("expected %q, got %q", c.err, err)
			}
			if c.output != w.String() {
				t.Errorf("expected %q, got %q", c.output, w.String())
			}
			if c.expectedName != name {
				t.Errorf("expected %q, got %q", c.input, name)
			}
		})
	}
}

type TimedReader struct {
	duration time.Duration
}

func (tr *TimedReader) Read(p []byte) (n int, err error) {
	time.Sleep(tr.duration)
	return 0, nil
}
