package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/shahzaibaziz/active-cookies/config"
	"github.com/shahzaibaziz/active-cookies/internal/app"
)

func main() {
	if err := run(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

// run contains the actual program logic, taking explicit writers instead of
// touching os.Stdout/os.Stderr directly so it can be exercised in tests
// without capturing global state.
func run(args []string, stdout, stderr io.Writer) error {
	cfg, err := config.Parse(args)
	if err != nil {
		return err
	}

	file, err := cfg.OpenFile()
	if err != nil {
		return err
	}
	defer file.Close()

	logger := log.New(stderr, "", 0)

	cookies, parseErrs := app.FindMostActive(file, cfg.Date)
	for _, e := range parseErrs {
		logger.Printf("skipping malformed line: %v", e)
	}

	for _, cookie := range cookies {
		if _, err := fmt.Fprintln(stdout, cookie); err != nil {
			return err
		}
	}

	return nil
}
