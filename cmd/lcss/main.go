package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"hash"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"lcss/internal/compile"
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
	case "build":
		if err := runBuild(os.Args[2:]); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	case "watch":
		if err := runWatch(os.Args[2:]); err != nil {
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

	sitePath := flags.String("site", "", "Path to site override config JSON (optional)")

	flags.Usage = func() {
		_, _ = fmt.Fprintln(os.Stdout, "Usage:")
		_, _ = fmt.Fprintln(os.Stdout, "  lcss config print [--site <path>]")
		_, _ = fmt.Fprintln(os.Stdout, "")
		_, _ = fmt.Fprintln(os.Stdout, "Options:")
		flags.PrintDefaults()
	}

	if err := flags.Parse(args); err != nil {
		return err
	}

	resolvedSitePath := resolveSitePath(*sitePath)
	cfg, err := config.Load("", resolvedSitePath)
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
	_, _ = fmt.Fprintln(os.Stdout, "  lcss config print [--site <path>]")
	_, _ = fmt.Fprintln(os.Stdout, "  lcss tokens [--site <path>] [--out <path>]")
	_, _ = fmt.Fprintln(os.Stdout, "  lcss build [--site <path>] [--out <path>] [--stdout] [--production]")
	_, _ = fmt.Fprintln(os.Stdout, "  lcss watch [--site <path>] [--out <path>] [--interval <dur>] [--once]")
	_, _ = fmt.Fprintln(os.Stdout, "  lcss scan [--site <path>] [--top <n>]")
}

func printConfigUsage() {
	_, _ = fmt.Fprintln(os.Stdout, "Usage:")
	_, _ = fmt.Fprintln(os.Stdout, "  lcss config print [--site <path>]")
}

func runTokens(args []string) error {
	flags := flag.NewFlagSet("tokens", flag.ContinueOnError)
	flags.SetOutput(os.Stdout)

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

	resolvedSitePath := resolveSitePath(*sitePath)
	cfg, err := config.Load("", resolvedSitePath)
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
	_, _ = fmt.Fprintln(os.Stdout, "  lcss tokens [--site <path>] [--out <path>]")
}

func runBuild(args []string) error {
	flags := flag.NewFlagSet("build", flag.ContinueOnError)
	flags.SetOutput(os.Stdout)

	sitePath := flags.String("site", "", "Path to site override config JSON (optional)")
	outPath := flags.String("out", "dist/lattice.css", "Path to output CSS")
	stdout := flags.Bool("stdout", false, "Write CSS to stdout instead of a file")
	production := flags.Bool("production", false, "Emit only classes found in build.content (production build)")

	flags.Usage = func() {
		printBuildUsage()
		_, _ = fmt.Fprintln(os.Stdout, "")
		_, _ = fmt.Fprintln(os.Stdout, "Options:")
		flags.PrintDefaults()
	}

	if err := flags.Parse(args); err != nil {
		return err
	}

	resolvedSitePath := resolveSitePath(*sitePath)
	cfg, err := config.Load("", resolvedSitePath)
	if err != nil {
		return err
	}
	if err := cfg.Validate(); err != nil {
		return err
	}

	canonical := cfg.Canonicalize()
	var output compile.Output
	if *production {
		if len(cfg.Build.Content) == 0 {
			return errors.New("build.content is required for --production")
		}
		result, err := extract.FromPaths(cfg.Build.Content, cfg.Build.Safelist)
		if err != nil {
			return err
		}
		output, err = compile.Build(canonical, result)
		if err != nil {
			return err
		}
	} else {
		classes := compile.AllClasses(canonical)
		result := extract.Result{Classes: classes}
		output, err = compile.Build(canonical, result)
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}
	for _, warning := range output.Warnings {
		_, _ = fmt.Fprintln(os.Stderr, "warning:", warning)
	}

	if *stdout {
		if len(output.Manifest) > 0 {
			_, _ = fmt.Fprintln(os.Stderr, "warning: manifest is not written when --stdout is set")
		}
		_, err := os.Stdout.Write(output.CSS)
		return err
	}

	return emit.Write(emit.Artifacts{LatticeCSS: output.CSS, Manifest: output.Manifest}, *outPath)
}

func printBuildUsage() {
	_, _ = fmt.Fprintln(os.Stdout, "Usage:")
	_, _ = fmt.Fprintln(os.Stdout, "  lcss build [--site <path>] [--out <path>] [--stdout] [--production]")
}

func runWatch(args []string) error {
	flags := flag.NewFlagSet("watch", flag.ContinueOnError)
	flags.SetOutput(os.Stdout)

	sitePath := flags.String("site", "", "Path to site override config JSON (optional)")
	outPath := flags.String("out", "dist/lattice.css", "Path to output CSS")
	interval := flags.Duration("interval", 500*time.Millisecond, "Polling interval for changes")
	once := flags.Bool("once", false, "Run one incremental check and exit")

	flags.Usage = func() {
		printWatchUsage()
		_, _ = fmt.Fprintln(os.Stdout, "")
		_, _ = fmt.Fprintln(os.Stdout, "Options:")
		flags.PrintDefaults()
	}

	if err := flags.Parse(args); err != nil {
		return err
	}
	if *interval <= 0 {
		return errors.New("--interval must be greater than zero")
	}

	state := watchState{}
	cachePath := cacheFilePath(*outPath)

	for {
		resolvedSitePath := resolveSitePath(*sitePath)
		built, err := watchOnce("", resolvedSitePath, *outPath, cachePath, &state)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "watch error:", err)
			if *once {
				return err
			}
		}
		if built {
			_, _ = fmt.Fprintln(os.Stdout, "built")
		}
		if *once {
			return nil
		}
		time.Sleep(*interval)
	}
}

func printWatchUsage() {
	_, _ = fmt.Fprintln(os.Stdout, "Usage:")
	_, _ = fmt.Fprintln(os.Stdout, "  lcss watch [--site <path>] [--out <path>] [--interval <dur>] [--once]")
}

func runScan(args []string) error {
	flags := flag.NewFlagSet("scan", flag.ContinueOnError)
	flags.SetOutput(os.Stdout)

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

	resolvedSitePath := resolveSitePath(*sitePath)
	cfg, err := config.Load("", resolvedSitePath)
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
	_, _ = fmt.Fprintln(os.Stdout, "  lcss scan [--site <path>] [--top <n>] [--per-file]")
}

type classCount struct {
	Class string
	Count int
}

type watchState struct {
	lastHash string
}

type cacheRecord struct {
	Version   int    `json:"version"`
	InputHash string `json:"inputHash"`
}

const cacheVersion = 1
const cacheFileName = ".lcss.cache.json"

func watchOnce(basePath, sitePath, outPath, cachePath string, state *watchState) (bool, error) {
	cfg, err := config.Load(basePath, sitePath)
	if err != nil {
		return false, err
	}
	if err := cfg.Validate(); err != nil {
		return false, err
	}
	if len(cfg.Build.Content) == 0 {
		return false, errors.New("build.content is required")
	}

	files, err := extract.FilesFromPatterns(cfg.Build.Content)
	if err != nil {
		return false, err
	}
	sort.Strings(files)

	hash, err := computeInputHash(basePath, sitePath, files)
	if err != nil {
		return false, err
	}

	if state.lastHash == "" {
		cache, err := loadCache(cachePath)
		if err != nil {
			return false, err
		}
		if cacheValid(cache, hash, outPath, cfg.Build.Emit.Manifest) {
			state.lastHash = hash
			return false, nil
		}
	}

	if hash == state.lastHash {
		return false, nil
	}

	result, err := extract.FromPaths(cfg.Build.Content, cfg.Build.Safelist)
	if err != nil {
		return false, err
	}

	canonical := cfg.Canonicalize()
	output, err := compile.Build(canonical, result)
	if err != nil {
		return false, err
	}
	for _, warning := range output.Warnings {
		_, _ = fmt.Fprintln(os.Stderr, "warning:", warning)
	}

	if err := emit.Write(emit.Artifacts{LatticeCSS: output.CSS, Manifest: output.Manifest}, outPath); err != nil {
		return false, err
	}
	if err := writeCache(cachePath, hash); err != nil {
		return false, err
	}

	state.lastHash = hash
	return true, nil
}

func computeInputHash(basePath, sitePath string, contentFiles []string) (string, error) {
	hasher := sha256.New()
	if basePath == "" {
		if err := hashBytesContents(hasher, "embedded:default", config.DefaultJSON()); err != nil {
			return "", err
		}
	} else {
		if err := hashFileContents(hasher, basePath); err != nil {
			return "", err
		}
	}
	if sitePath != "" {
		if err := hashFileContents(hasher, sitePath); err != nil {
			return "", err
		}
	}

	for _, path := range contentFiles {
		info, err := os.Stat(path)
		if err != nil {
			return "", err
		}
		_, _ = fmt.Fprintf(hasher, "file:%s\n", path)
		_, _ = fmt.Fprintf(hasher, "size:%d\n", info.Size())
		_, _ = fmt.Fprintf(hasher, "mtime:%d\n", info.ModTime().UnixNano())
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func hashFileContents(hasher hash.Hash, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	_, _ = fmt.Fprintf(hasher, "config:%s\n", path)
	_, _ = hasher.Write(data)
	_, _ = fmt.Fprintln(hasher)
	return nil
}

func hashBytesContents(hasher hash.Hash, label string, data []byte) error {
	_, _ = fmt.Fprintf(hasher, "config:%s\n", label)
	_, _ = hasher.Write(data)
	_, _ = fmt.Fprintln(hasher)
	return nil
}

func resolveSitePath(path string) string {
	if path != "" {
		return path
	}
	if fileExists("lattice.json") {
		return "lattice.json"
	}
	return ""
}

func cacheFilePath(outPath string) string {
	return filepath.Join(filepath.Dir(outPath), cacheFileName)
}

func loadCache(path string) (cacheRecord, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return cacheRecord{}, nil
		}
		return cacheRecord{}, err
	}
	var record cacheRecord
	if err := json.Unmarshal(data, &record); err != nil {
		return cacheRecord{}, err
	}
	return record, nil
}

func writeCache(path, hash string) error {
	record := cacheRecord{Version: cacheVersion, InputHash: hash}
	data, err := json.Marshal(record)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func cacheValid(record cacheRecord, hash, outPath string, manifestExpected bool) bool {
	if record.Version != cacheVersion || record.InputHash == "" {
		return false
	}
	if record.InputHash != hash {
		return false
	}
	return outputsExist(outPath, manifestExpected)
}

func outputsExist(outPath string, manifestExpected bool) bool {
	if !fileExists(outPath) {
		return false
	}
	if !manifestExpected {
		return true
	}
	manifestPath := filepath.Join(filepath.Dir(outPath), "manifest.json")
	return fileExists(manifestPath)
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}
