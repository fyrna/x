package sh

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Options struct {
	Env      []string
	Silent   bool
	Capture  bool   // true = capture command output
	Combined bool   // true = capture both stdout and stderr together
	Shell    string // override shell (e.g. "bash", "sh", "zsh")
}

// go dont have "pub" keyword, and basically it is illegal to use dollar sign. :(
// i was want something like bun does $`echo rawrrr`
var S = func(cmd string) error {
	_, err := Exec(context.Background(), &Options{}, cmd)
	return err
}

// Exec - main function that handles all use cases
func Exec(ctx context.Context, opts *Options, command string) (string, error) {
	if opts == nil {
		opts = &Options{}
	}

	// detect if shell is needed (redirection, pipes, etc.)
	needsShell := strings.ContainsAny(command, "><|&;{}*?[]$`()")

	var cmd *exec.Cmd
	var err error

	if needsShell || opts.Shell != "" {
		shell := opts.Shell

		// default: try bash, fallback to sh
		if shell == "" {
			if _, err := exec.LookPath("bash"); err == nil {
				shell = "bash"
			} else {
				shell = "sh"
			}
		}

		cmd = exec.CommandContext(ctx, shell, "-c", command)
	} else {
		// direct exec (no shell involved)
		parts, err := splitCommand(command)
		if err != nil {
			return "", fmt.Errorf("execute: %w", err)
		}
		if len(parts) == 0 {
			return "", nil
		}
		cmd = exec.CommandContext(ctx, parts[0], parts[1:]...)
	}

	// set environment variables if provided
	if opts.Env != nil {
		cmd.Env = append(os.Environ(), opts.Env...)
	}

	var output []byte
	if opts.Capture {
		if opts.Combined {
			output, err = cmd.CombinedOutput()
		} else {
			output, err = cmd.Output()
		}
	} else {
		if !opts.Silent {
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
		}
		err = cmd.Run()
	}

	if err != nil {
		return "", fmt.Errorf("execute %s: %w", command, err)
	}
	return string(output), nil
}

// Sh - convenience wrapper for a raw string command
func Sh(ctx context.Context, cmd string) error {
	_, err := Exec(ctx, &Options{}, cmd)
	return err
}

// ShEnv - convenience wrapper for a raw string command with env vars
func ShEnv(ctx context.Context, env []string, cmd string) error {
	_, err := Exec(ctx, &Options{Env: env}, cmd)
	return err
}

// ShOut - convenience wrapper for a raw string command that returns output
func Output(ctx context.Context, cmd string) (string, error) {
	return Exec(ctx, &Options{Capture: true}, cmd)
}

// Silent - run a command with env vars, no output printed
func Silent(ctx context.Context, env []string, cmd string, args ...string) error {
	_, err := Exec(ctx, &Options{
		Env:    env,
		Silent: true,
	}, cmd)
	return err
}

func splitCommand(input string) ([]string, error) {
	var args []string
	var current []rune

	inSingle := false
	inDouble := false
	escape := false

	for _, c := range input {
		switch {
		case escape:
			current = append(current, c)
			escape = false
		case c == '\\':
			escape = true
		case c == '"' && !inSingle:
			inDouble = !inDouble
		case c == '\'' && !inDouble:
			inSingle = !inSingle
		case c == ' ' && !inSingle && !inDouble:
			if len(current) > 0 {
				args = append(args, string(current))
				current = nil
			}
		default:
			current = append(current, c)
		}
	}

	if len(current) > 0 {
		args = append(args, string(current))
	}

	if inSingle || inDouble {
		return nil, fmt.Errorf("unmatched quotes in command: %s", input)
	}

	if escape {
		return nil, fmt.Errorf("unfinished escape sequence in command: %s", input)
	}

	return args, nil
}
