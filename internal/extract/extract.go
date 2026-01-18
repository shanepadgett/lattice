package extract

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var (
	classAttrPattern     = regexp.MustCompile(`(?s)\bclass\s*=\s*(?:"([^"]*)"|'([^']*)')`)
	classNameAttrPattern = regexp.MustCompile(`(?s)\bclassName\s*=\s*(?:"([^"]*)"|'([^']*)')`)
	classActionPattern   = regexp.MustCompile(`(?s)\bclass\s*=\s*{{(.*?)}}`)
	classNameAction      = regexp.MustCompile(`(?s)\bclassName\s*=\s*{{(.*?)}}`)
	stringLiteralPattern = regexp.MustCompile(`"(?:\\.|[^"\\])*"|` + "`" + `[^` + "`" + `]*` + "`")
	validClassPattern    = regexp.MustCompile(`^[a-zA-Z0-9-:_]+$`)
)

type Result struct {
	Classes []string
	Counts  map[string]int
	ByFile  map[string]map[string]int
	Files   int
}

func FromPaths(patterns []string, safelist []string) (Result, error) {
	files, err := expandPatterns(patterns)
	if err != nil {
		return Result{}, err
	}

	classSet := map[string]struct{}{}
	counts := map[string]int{}
	byFile := map[string]map[string]int{}

	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			return Result{}, fmt.Errorf("read %s: %w", file, err)
		}

		fileCounts := map[string]int{}

		for _, class := range extractClasses(string(data)) {
			if !validClassPattern.MatchString(class) {
				continue
			}
			counts[class]++
			classSet[class] = struct{}{}
			fileCounts[class]++
		}

		if len(fileCounts) > 0 {
			byFile[file] = fileCounts
		}
	}

	for _, class := range safelist {
		class = strings.TrimSpace(class)
		if class == "" || !validClassPattern.MatchString(class) {
			continue
		}
		counts[class]++
		classSet[class] = struct{}{}
	}

	classes := make([]string, 0, len(classSet))
	for class := range classSet {
		classes = append(classes, class)
	}
	sort.Strings(classes)

	return Result{
		Classes: classes,
		Counts:  counts,
		ByFile:  byFile,
		Files:   len(files),
	}, nil
}

func expandPatterns(patterns []string) ([]string, error) {
	seen := map[string]struct{}{}
	var files []string

	for _, pattern := range patterns {
		pattern = strings.TrimSpace(pattern)
		if pattern == "" {
			continue
		}

		matches, err := globPattern(pattern)
		if err != nil {
			return nil, fmt.Errorf("glob %s: %w", pattern, err)
		}
		for _, match := range matches {
			match = filepath.Clean(match)
			if _, ok := seen[match]; ok {
				continue
			}
			info, err := os.Stat(match)
			if err != nil {
				continue
			}
			if info.IsDir() {
				continue
			}
			seen[match] = struct{}{}
			files = append(files, match)
		}
	}

	return files, nil
}

func globPattern(pattern string) ([]string, error) {
	if !strings.Contains(pattern, "**") {
		return filepath.Glob(pattern)
	}

	root, remaining := splitGlob(pattern)
	if root == "" {
		root = "."
	}

	rootInfo, err := os.Stat(root)
	if err != nil {
		return nil, err
	}
	if !rootInfo.IsDir() {
		return nil, nil
	}

	var matches []string
	walkErr := filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		if globMatch(filepath.ToSlash(remaining), filepath.ToSlash(rel)) {
			matches = append(matches, path)
		}
		return nil
	})
	if walkErr != nil {
		return nil, walkErr
	}
	return matches, nil
}

func splitGlob(pattern string) (string, string) {
	pattern = filepath.Clean(pattern)
	patternSlash := filepath.ToSlash(pattern)
	segments := strings.Split(patternSlash, "/")

	rootEnd := 0
	for _, segment := range segments {
		if strings.ContainsAny(segment, "*?[]") {
			break
		}
		rootEnd++
	}

	if rootEnd == 0 {
		return ".", pattern
	}

	root := filepath.FromSlash(strings.Join(segments[:rootEnd], "/"))
	remaining := strings.Join(segments[rootEnd:], "/")
	if remaining == "" {
		remaining = "**"
	}
	return root, remaining
}

func globMatch(pattern, path string) bool {
	pattern = strings.TrimPrefix(pattern, "./")
	path = strings.TrimPrefix(path, "./")

	p := splitGlobSegments(pattern)
	s := splitGlobSegments(path)

	return matchSegments(p, s)
}

func splitGlobSegments(value string) []string {
	if value == "" {
		return nil
	}
	return strings.Split(value, "/")
}

func matchSegments(pattern []string, path []string) bool {
	for len(pattern) > 0 {
		segment := pattern[0]
		pattern = pattern[1:]

		if segment == "**" {
			if len(pattern) == 0 {
				return true
			}
			for i := 0; i <= len(path); i++ {
				if matchSegments(pattern, path[i:]) {
					return true
				}
			}
			return false
		}

		if len(path) == 0 {
			return false
		}

		matched, err := filepath.Match(segment, path[0])
		if err != nil || !matched {
			return false
		}
		path = path[1:]
	}

	return len(path) == 0
}

func extractClasses(content string) []string {
	var classes []string

	classes = append(classes, extractAttrClasses(classAttrPattern, content)...)
	classes = append(classes, extractAttrClasses(classNameAttrPattern, content)...)
	classes = append(classes, extractActionClasses(classActionPattern, content)...)
	classes = append(classes, extractActionClasses(classNameAction, content)...)

	return classes
}

func extractAttrClasses(pattern *regexp.Regexp, content string) []string {
	var classes []string
	for _, match := range pattern.FindAllStringSubmatch(content, -1) {
		value := match[1]
		if value == "" {
			value = match[2]
		}
		classes = append(classes, splitClasses(value)...)
	}
	return classes
}

func extractActionClasses(pattern *regexp.Regexp, content string) []string {
	var classes []string
	for _, match := range pattern.FindAllStringSubmatch(content, -1) {
		for _, literal := range stringLiteralPattern.FindAllString(match[1], -1) {
			classes = append(classes, splitClasses(stripQuotes(literal))...)
		}
	}
	return classes
}

func stripQuotes(value string) string {
	if len(value) < 2 {
		return value
	}
	return value[1 : len(value)-1]
}

func splitClasses(value string) []string {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	return strings.Fields(value)
}
