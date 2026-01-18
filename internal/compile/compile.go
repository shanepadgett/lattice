package compile

import (
	"fmt"
	"strings"

	"lcss/internal/config"
	"lcss/internal/emit"
	"lcss/internal/extract"
)

const manifestVersion = 1

type Output struct {
	CSS      []byte
	Manifest []byte
	Warnings []string
}

func Build(canonical config.Canonical, result extract.Result) (Output, error) {
	policy := canonical.Config.Build.UnknownClassPolicy
	if policy == "" {
		policy = "warn"
	}

	sections := make([]string, 0, 3)
	if canonical.Config.Build.Emit.TokensCSS {
		tokens, err := emit.TokensCSS(canonical)
		if err != nil {
			return Output{}, fmt.Errorf("emit tokens: %w", err)
		}
		sections = append(sections, strings.TrimRight(string(tokens), "\n"))
	}

	if canonical.Config.Build.Emit.Base {
		sections = append(sections, strings.TrimRight(baseCSS(canonical), "\n"))
	}

	rules, matched, unknown := buildUtilities(canonical, result.Classes)
	utilities := strings.TrimRight(renderRules(rules), "\n")
	if utilities != "" {
		sections = append(sections, utilities)
	}

	var css string
	if len(sections) > 0 {
		css = strings.Join(sections, "\n\n") + "\n"
	}

	warnings := []string{}
	if policy == "warn" {
		for _, class := range unknown {
			warnings = append(warnings, fmt.Sprintf("unknown class: %s", class))
		}
	}
	if policy == "error" && len(unknown) > 0 {
		return Output{}, fmt.Errorf("unknown classes: %s", strings.Join(unknown, ", "))
	}

	var manifest []byte
	if canonical.Config.Build.Emit.Manifest {
		data, err := buildManifest(result, matched, unknown)
		if err != nil {
			return Output{}, fmt.Errorf("build manifest: %w", err)
		}
		manifest = data
	}

	return Output{
		CSS:      []byte(css),
		Manifest: manifest,
		Warnings: warnings,
	}, nil
}
