package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"lcss/internal/config"
	"lcss/internal/emit"
)

var version = "dev"
var commit = ""
var date = ""

func main() {
	if len(os.Args) == 1 || strings.HasPrefix(os.Args[1], "-") {
		runRoot(os.Args[1:])
		return
	}

	switch os.Args[1] {
	case "config":
		if err := runConfig(os.Args[2:]); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case "tokens":
		if err := runTokens(os.Args[2:]); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	default:
		printRootUsage()
		os.Exit(1)
	}
}

func runRoot(args []string) {
	flags := flag.NewFlagSet("lcss", flag.ContinueOnError)
	flags.SetOutput(os.Stdout)

	versionFlag := flags.Bool("version", false, "Print version and exit")
	versionJSON := flags.Bool("version-json", false, "Print version info as JSON and exit")

	flags.Usage = func() {
		printRootUsage()
		_, _ = fmt.Fprintln(os.Stdout, "")
		_, _ = fmt.Fprintln(os.Stdout, "Options:")
		flags.PrintDefaults()
	}

	if err := flags.Parse(args); err != nil {
		return
	}

	if *versionFlag || *versionJSON {
		info := struct {
			Version string `json:"version"`
			Commit  string `json:"commit"`
			Date    string `json:"date"`
		}{
			Version: version,
			Commit:  commit,
			Date:    date,
		}

		if *versionJSON {
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "  ")
			_ = enc.Encode(info)
			os.Exit(0)
		}

		fmt.Printf("lcss %s\ncommit: %s\nbuilt: %s\n", info.Version, info.Commit, info.Date)
		os.Exit(0)
	}

	flags.Usage()
}

func runConfig(args []string) error {
	if len(args) == 0 {
		printConfigUsage()
		return errors.New("config command requires a subcommand")
	}

	switch args[0] {
	case "print":
		return runConfigPrint(args[1:])
	default:
		printConfigUsage()
		return fmt.Errorf("unknown config subcommand: %s", args[0])
	}
}

func runConfigPrint(args []string) error {
	flags := flag.NewFlagSet("config print", flag.ContinueOnError)
	flags.SetOutput(os.Stdout)

	basePath := flags.String("base", "", "Path to base config JSON")
	sitePath := flags.String("site", "", "Path to site override config JSON (optional)")

	flags.Usage = func() {
		_, _ = fmt.Fprintln(os.Stdout, "Usage:")
		_, _ = fmt.Fprintln(os.Stdout, "  lcss config print --base <path> [--site <path>]")
		_, _ = fmt.Fprintln(os.Stdout, "")
		_, _ = fmt.Fprintln(os.Stdout, "Options:")
		flags.PrintDefaults()
	}

	if err := flags.Parse(args); err != nil {
		return err
	}
	if *basePath == "" {
		flags.Usage()
		return errors.New("--base is required")
	}

	cfg, err := config.Load(*basePath, *sitePath)
	if err != nil {
		return err
	}
	if err := cfg.Validate(); err != nil {
		return err
	}

	canonical := cfg.Canonicalize()
	data, err := config.MarshalDeterministic(canonical)
	if err != nil {
		return err
	}
	_, err = os.Stdout.Write(data)
	return err
}

func printRootUsage() {
	_, _ = fmt.Fprintln(os.Stdout, "lcss - Lattice CSS compiler")
	_, _ = fmt.Fprintln(os.Stdout, "")
	_, _ = fmt.Fprintln(os.Stdout, "Usage:")
	_, _ = fmt.Fprintln(os.Stdout, "  lcss [options]")
	_, _ = fmt.Fprintln(os.Stdout, "  lcss config print --base <path> [--site <path>]")
	_, _ = fmt.Fprintln(os.Stdout, "  lcss tokens --base <path> [--site <path>] [--out <path>]")
}

func printConfigUsage() {
	_, _ = fmt.Fprintln(os.Stdout, "Usage:")
	_, _ = fmt.Fprintln(os.Stdout, "  lcss config print --base <path> [--site <path>]")
}

func runTokens(args []string) error {
	flags := flag.NewFlagSet("tokens", flag.ContinueOnError)
	flags.SetOutput(os.Stdout)

	basePath := flags.String("base", "", "Path to base config JSON")
	sitePath := flags.String("site", "", "Path to site override config JSON (optional)")
	outPath := flags.String("out", "dist/tokens.css", "Path to output tokens CSS")
	stdout := flags.Bool("stdout", false, "Write tokens CSS to stdout instead of a file")

	flags.Usage = func() {
		printTokensUsage()
		_, _ = fmt.Fprintln(os.Stdout, "")
		_, _ = fmt.Fprintln(os.Stdout, "Options:")
		flags.PrintDefaults()
	}

	if err := flags.Parse(args); err != nil {
		return err
	}
	if *basePath == "" {
		flags.Usage()
		return errors.New("--base is required")
	}

	cfg, err := config.Load(*basePath, *sitePath)
	if err != nil {
		return err
	}
	if err := cfg.Validate(); err != nil {
		return err
	}

	canonical := cfg.Canonicalize()
	data, err := emit.TokensCSS(canonical)
	if err != nil {
		return err
	}

	if *stdout {
		_, err := os.Stdout.Write(data)
		return err
	}

	if err := os.MkdirAll(filepath.Dir(*outPath), 0o755); err != nil {
		return err
	}
	return os.WriteFile(*outPath, data, 0o644)
}

func printTokensUsage() {
	_, _ = fmt.Fprintln(os.Stdout, "Usage:")
	_, _ = fmt.Fprintln(os.Stdout, "  lcss tokens --base <path> [--site <path>] [--out <path>]")
}
