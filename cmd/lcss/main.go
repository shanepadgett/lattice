package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"lcss/internal/config"
	"lcss/internal/emit"
	"lcss/internal/extract"
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
	case "scan":
		if err := runScan(os.Args[2:]); err != nil {
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
	_, _ = fmt.Fprintln(os.Stdout, "  lcss scan --base <path> [--site <path>] [--top <n>]")
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

func runScan(args []string) error {
	flags := flag.NewFlagSet("scan", flag.ContinueOnError)
	flags.SetOutput(os.Stdout)

	basePath := flags.String("base", "", "Path to base config JSON")
	sitePath := flags.String("site", "", "Path to site override config JSON (optional)")
	top := flags.Int("top", 20, "Show top N classes by frequency (0 to disable)")
	perFile := flags.Bool("per-file", false, "Show top N classes per file")

	flags.Usage = func() {
		printScanUsage()
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
	if len(cfg.Build.Content) == 0 {
		return errors.New("build.content is required for scan")
	}

	result, err := extract.FromPaths(cfg.Build.Content, cfg.Build.Safelist)
	if err != nil {
		return err
	}

	_, _ = fmt.Fprintf(os.Stdout, "Files: %d\n", result.Files)
	_, _ = fmt.Fprintf(os.Stdout, "Classes: %d\n", len(result.Classes))
	if *top == 0 {
		return nil
	}

	items := make([]classCount, 0, len(result.Counts))
	for class, count := range result.Counts {
		items = append(items, classCount{Class: class, Count: count})
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].Count == items[j].Count {
			return items[i].Class < items[j].Class
		}
		return items[i].Count > items[j].Count
	})

	limit := *top
	if limit > len(items) {
		limit = len(items)
	}
	if limit > 0 {
		_, _ = fmt.Fprintf(os.Stdout, "Top %d:\n", limit)
		for _, item := range items[:limit] {
			_, _ = fmt.Fprintf(os.Stdout, "  %s (%d)\n", item.Class, item.Count)
		}
	}

	if !*perFile {
		return nil
	}

	fileNames := make([]string, 0, len(result.ByFile))
	for name := range result.ByFile {
		fileNames = append(fileNames, name)
	}
	sort.Strings(fileNames)

	for _, name := range fileNames {
		perFileItems := make([]classCount, 0, len(result.ByFile[name]))
		for class, count := range result.ByFile[name] {
			perFileItems = append(perFileItems, classCount{Class: class, Count: count})
		}
		sort.Slice(perFileItems, func(i, j int) bool {
			if perFileItems[i].Count == perFileItems[j].Count {
				return perFileItems[i].Class < perFileItems[j].Class
			}
			return perFileItems[i].Count > perFileItems[j].Count
		})

		perFileLimit := *top
		if perFileLimit > len(perFileItems) {
			perFileLimit = len(perFileItems)
		}
		if perFileLimit == 0 {
			continue
		}

		_, _ = fmt.Fprintf(os.Stdout, "%s:\n", name)
		for _, item := range perFileItems[:perFileLimit] {
			_, _ = fmt.Fprintf(os.Stdout, "  %s (%d)\n", item.Class, item.Count)
		}
	}
	return nil
}

func printScanUsage() {
	_, _ = fmt.Fprintln(os.Stdout, "Usage:")
	_, _ = fmt.Fprintln(os.Stdout, "  lcss scan --base <path> [--site <path>] [--top <n>] [--per-file]")
}

type classCount struct {
	Class string
	Count int
}
