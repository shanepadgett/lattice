package emit

import (
	"sort"
	"strings"

	"lcss/internal/config"
)

type tokenEntry struct {
	name  string
	value string
}

func TokensCSS(canonical config.Canonical) ([]byte, error) {
	entries := make([]tokenEntry, 0)
	entries = appendScaleTokens(entries, canonical.Tokens.Scales)
	entries = appendThemeTokens(entries, canonical.Tokens.Themes["default"])

	var b strings.Builder
	writeBlock(&b, ":root", entries)

	themeNames := make([]string, 0, len(canonical.Tokens.Themes))
	for name := range canonical.Tokens.Themes {
		if name == "default" {
			continue
		}
		themeNames = append(themeNames, name)
	}
	sort.Strings(themeNames)

	for _, name := range themeNames {
		themeEntries := appendThemeTokens(nil, canonical.Tokens.Themes[name])
		if len(themeEntries) == 0 {
			continue
		}
		b.WriteString("\n")
		selector := "[data-theme=\"" + name + "\"]"
		writeBlock(&b, selector, themeEntries)
	}

	return []byte(b.String()), nil
}

func writeBlock(b *strings.Builder, selector string, entries []tokenEntry) {
	b.WriteString(selector)
	b.WriteString(" {\n")
	for _, entry := range entries {
		b.WriteString("  ")
		b.WriteString(entry.name)
		b.WriteString(": ")
		b.WriteString(entry.value)
		b.WriteString(";\n")
	}
	b.WriteString("}\n")
}

func appendScaleTokens(entries []tokenEntry, scales map[string]map[string]string) []tokenEntry {
	scaleOrder := []string{
		"space",
		"size",
		"radius",
		"borderWidth",
		"fontSize",
		"lineHeight",
		"fontWeight",
		"letterSpacing",
		"shadow",
		"z",
		"opacity",
		"aspect",
		"duration",
		"easing",
		"delay",
		"translate",
		"rotate",
		"scale",
		"maxWidth",
		"maxHeight",
		"container",
	}
	for _, scale := range scaleOrder {
		values, ok := scales[scale]
		if !ok || len(values) == 0 {
			continue
		}
		entries = appendTokenMap(entries, scalePrefix(scale), values)
	}
	return entries
}

func appendThemeTokens(entries []tokenEntry, theme config.ThemeTokens) []tokenEntry {
	if len(theme.Colors) > 0 {
		entries = appendTokenMap(entries, "color", theme.Colors)
	}
	if len(theme.Fonts) > 0 {
		entries = appendTokenMap(entries, "font", theme.Fonts)
	}
	return entries
}

func appendTokenMap(entries []tokenEntry, prefix string, values map[string]string) []tokenEntry {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		entries = append(entries, tokenEntry{
			name:  "--" + prefix + "-" + key,
			value: values[key],
		})
	}
	return entries
}

func scalePrefix(scale string) string {
	switch scale {
	case "fontSize":
		return "font-size"
	case "lineHeight":
		return "line-height"
	case "fontWeight":
		return "font-weight"
	case "borderWidth":
		return "border-width"
	case "letterSpacing":
		return "letter-spacing"
	case "maxWidth":
		return "max-width"
	case "maxHeight":
		return "max-height"
	default:
		return scale
	}
}
