package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/shota3506/gostlc/internal/eval"
	"github.com/shota3506/gostlc/internal/parser"
	"github.com/shota3506/gostlc/internal/types"
)

var (
	command = flag.String("c", "", "Execute STLC code from command line")
	help    = flag.Bool("h", false, "Show help")
)

func main() {
	flag.Parse()

	if *help {
		usage()
		os.Exit(0)
	}

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	if *command != "" {
		return runCode(*command)
	}

	args := flag.Args()

	switch len(args) {
	case 0:
		if isInteractive() {
			return startREPL()
		} else {
			return runStdin()
		}
	case 1:
		if args[0] == "-" {
			return runStdin()
		} else {
			return runFile(args[0])
		}
	default:
		return errors.New("too many arguments")
	}
}

func usage() {
	command := "gostlc"
	fmt.Fprintf(os.Stderr, "Usage: %s [options] [file]\n", command)
	fmt.Fprintf(os.Stderr, "\nOptions:\n")
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\nExamples:\n")
	fmt.Fprintf(os.Stderr, "  %s                    # Start REPL\n", command)
	fmt.Fprintf(os.Stderr, "  %s file.stlc          # Run file\n", command)
	fmt.Fprintf(os.Stderr, "  %s -c \"(\\x:Int.x) 42\" # Execute code\n", command)
	fmt.Fprintf(os.Stderr, "  echo \"code\" | %s -    # Read from stdin\n", command)
}

func isInteractive() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

func runFile(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return runCode(string(data))
}

func runStdin() error {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		return err
	}
	return runCode(string(data))
}

func runCode(code string) error {
	expr, err := parser.Parse(code)
	if err != nil {
		return err
	}

	typedExpr, err := types.Check(expr)
	if err != nil {
		return err
	}

	resp, err := eval.Eval(typedExpr)
	if err != nil {
		return err
	}

	fmt.Fprintln(os.Stdout, resp.String())
	return nil
}

func startREPL() error {
	fmt.Println("STLC REPL (type :quit to exit)")
	fmt.Println("REPL not yet implemented")
	// TODO: Implement REPL
	return nil
}
