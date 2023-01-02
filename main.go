package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

var totalDuration time.Duration = 5

func main() {
	allowedDuration := totalDuration * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), allowedDuration)
	defer cancel()

	name, err := getNameContext(ctx, os.Stdin, os.Stdout)
	if err != nil && !errors.Is(err, context.DeadlineExceeded) {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	fmt.Fprintln(os.Stdout, fmt.Sprintf("Hello %s\n", name))
}

func getNameContext(ctx context.Context, r io.Reader, w io.Writer) (string, error) {
	var err error
	name := "Default Name"
	c := make(chan error, 1)

	go func() {
		name, err = getName(r, w)
		c <- err
	}()

	select {
	case <-ctx.Done():
		return name, ctx.Err()
	case <-c:
		return name, err
	}
}

func getName(r io.Reader, w io.Writer) (string, error) {
	scanner := bufio.NewScanner(r)
	msg := "Your name please? Press the key Enter when done"
	fmt.Fprintln(w, msg)
	scanner.Scan()

	if err := scanner.Err(); err != nil {
		return "", err
	}

	name := scanner.Text()
	if len(name) == 0 {
		return "", errors.New("you entered an empty name")
	}

	return name, nil
}
