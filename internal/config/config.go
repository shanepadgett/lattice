package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
)

type Config struct {
	SchemaVersion int               `json:"schemaVersion"`
	ClassPrefix   string            `json:"classPrefix,omitempty"`
	Separator     string            `json:"separator,omitempty"`
	Breakpoints   map[string]string `json:"breakpoints,omitempty"`
	Themes        map[string]Theme  `json:"themes,omitempty"`
	Scales        Scales            `json:"scales,omitempty"`
	Variants      Variants          `json:"variants,omitempty"`
	Build         Build             `json:"build,omitempty"`
}

type Theme struct {
	Colors map[string]string `json:"colors,omitempty"`
	Font   map[string]string `json:"font,omitempty"`
}

type Scales struct {
	Space      map[string]string `json:"space,omitempty"`
	Radius     map[string]string `json:"radius,omitempty"`
	FontSize   map[string]string `json:"fontSize,omitempty"`
	LineHeight map[string]string `json:"lineHeight,omitempty"`
	FontWeight map[string]string `json:"fontWeight,omitempty"`
}

type Variants struct {
	Responsive []string `json:"responsive,omitempty"`
	State      []string `json:"state,omitempty"`
}

type Build struct {
	Content            []string    `json:"content,omitempty"`
	Safelist           []string    `json:"safelist,omitempty"`
	Emit               EmitOptions `json:"emit,omitempty"`
	UnknownClassPolicy string      `json:"unknownClassPolicy,omitempty"`
}

type EmitOptions struct {
	TokensCSS bool `json:"tokensCss,omitempty"`
	Base      bool `json:"base,omitempty"`
	Manifest  bool `json:"manifest,omitempty"`
}

type Canonical struct {
	Config
	Tokens CanonicalTokens `json:"tokens"`
}

type CanonicalTokens struct {
	Themes map[string]ThemeTokens       `json:"themes"`
	Scales map[string]map[string]string `json:"scales"`
}

type ThemeTokens struct {
	Colors map[string]string `json:"colors"`
	Fonts  map[string]string `json:"fonts"`
}

func Default() Config {
	return Config{
		Separator: ":",
	}
}

func Load(basePath, sitePath string) (Config, error) {
	if basePath == "" {
		return Config{}, errors.New("base config path is required")
	}

	baseMap, err := readJSONFile(basePath)
	if err != nil {
		return Config{}, fmt.Errorf("read base config: %w", err)
	}

	merged := baseMap
	if sitePath != "" {
		siteMap, err := readJSONFile(sitePath)
		if err != nil {
			return Config{}, fmt.Errorf("read site config: %w", err)
		}
		merged = mergeMaps(baseMap, siteMap)
	}

	data, err := json.Marshal(merged)
	if err != nil {
		return Config{}, fmt.Errorf("marshal merged config: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, fmt.Errorf("decode merged config: %w", err)
	}

	if cfg.Separator == "" {
		cfg.Separator = ":"
	}
	if cfg.Build.UnknownClassPolicy == "" {
		cfg.Build.UnknownClassPolicy = "warn"
	}

	return cfg, nil
}

func (c Config) Validate() error {
	if c.SchemaVersion == 0 {
		return errors.New("schemaVersion is required")
	}
	if c.Themes == nil || c.Themes["default"].Colors == nil && c.Themes["default"].Font == nil {
		if _, ok := c.Themes["default"]; !ok {
			return errors.New("themes.default is required")
		}
	}
	if len(c.Scales.Space) == 0 {
		return errors.New("scales.space is required")
	}
	if c.Build.UnknownClassPolicy != "" {
		switch c.Build.UnknownClassPolicy {
		case "ignore", "warn", "error":
			// ok
		default:
			return fmt.Errorf("build.unknownClassPolicy must be one of ignore, warn, error")
		}
	}
	if len(c.Variants.Responsive) > 0 {
		if c.Breakpoints == nil {
			return errors.New("breakpoints are required when variants.responsive is set")
		}
		for _, name := range c.Variants.Responsive {
			if _, ok := c.Breakpoints[name]; !ok {
				return fmt.Errorf("variants.responsive references unknown breakpoint: %s", name)
			}
		}
	}
	return nil
}

func (c Config) Canonicalize() Canonical {
	return Canonical{
		Config: c,
		Tokens: NormalizeTokens(c),
	}
}

func NormalizeTokens(c Config) CanonicalTokens {
	tokens := CanonicalTokens{
		Themes: map[string]ThemeTokens{},
		Scales: map[string]map[string]string{},
	}

	for name, theme := range c.Themes {
		tokens.Themes[name] = ThemeTokens{
			Colors: copyStringMap(theme.Colors),
			Fonts:  copyStringMap(theme.Font),
		}
	}

	if c.Scales.Space != nil {
		tokens.Scales["space"] = copyStringMap(c.Scales.Space)
	}
	if c.Scales.Radius != nil {
		tokens.Scales["radius"] = copyStringMap(c.Scales.Radius)
	}
	if c.Scales.FontSize != nil {
		tokens.Scales["fontSize"] = copyStringMap(c.Scales.FontSize)
	}
	if c.Scales.LineHeight != nil {
		tokens.Scales["lineHeight"] = copyStringMap(c.Scales.LineHeight)
	}
	if c.Scales.FontWeight != nil {
		tokens.Scales["fontWeight"] = copyStringMap(c.Scales.FontWeight)
	}

	return tokens
}

func MarshalDeterministic(v any) ([]byte, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	var decoded any
	if err := json.Unmarshal(data, &decoded); err != nil {
		return nil, err
	}

	buf := &bytes.Buffer{}
	if err := encodeSorted(buf, decoded, 0); err != nil {
		return nil, err
	}
	buf.WriteByte('\n')
	return buf.Bytes(), nil
}

func readJSONFile(path string) (map[string]any, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var out map[string]any
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, err
	}

	return out, nil
}

func mergeMaps(base, override map[string]any) map[string]any {
	if base == nil {
		base = map[string]any{}
	}
	if override == nil {
		return base
	}
	merged := make(map[string]any, len(base))
	for key, value := range base {
		merged[key] = value
	}
	for key, overrideValue := range override {
		if baseValue, ok := merged[key]; ok {
			merged[key] = mergeValue(baseValue, overrideValue)
			continue
		}
		merged[key] = overrideValue
	}
	return merged
}

func mergeValue(base, override any) any {
	if override == nil {
		return base
	}
	if overrideMap, ok := override.(map[string]any); ok {
		if baseMap, ok := base.(map[string]any); ok {
			return mergeMaps(baseMap, overrideMap)
		}
		return overrideMap
	}
	if overrideSlice, ok := override.([]any); ok {
		return overrideSlice
	}
	return override
}

func copyStringMap(input map[string]string) map[string]string {
	if input == nil {
		return nil
	}
	out := make(map[string]string, len(input))
	for key, value := range input {
		out[key] = value
	}
	return out
}

func encodeSorted(w io.Writer, value any, indent int) error {
	switch typed := value.(type) {
	case map[string]any:
		keys := make([]string, 0, len(typed))
		for key := range typed {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		if _, err := io.WriteString(w, "{"); err != nil {
			return err
		}
		if len(keys) > 0 {
			if _, err := io.WriteString(w, "\n"); err != nil {
				return err
			}
		}
		for i, key := range keys {
			writeIndent(w, indent+1)
			encodedKey, _ := json.Marshal(key)
			if _, err := w.Write(encodedKey); err != nil {
				return err
			}
			if _, err := io.WriteString(w, ": "); err != nil {
				return err
			}
			if err := encodeSorted(w, typed[key], indent+1); err != nil {
				return err
			}
			if i < len(keys)-1 {
				if _, err := io.WriteString(w, ","); err != nil {
					return err
				}
			}
			if _, err := io.WriteString(w, "\n"); err != nil {
				return err
			}
		}
		if len(keys) > 0 {
			writeIndent(w, indent)
		}
		_, err := io.WriteString(w, "}")
		return err
	case []any:
		if _, err := io.WriteString(w, "["); err != nil {
			return err
		}
		if len(typed) > 0 {
			if _, err := io.WriteString(w, "\n"); err != nil {
				return err
			}
		}
		for i, item := range typed {
			writeIndent(w, indent+1)
			if err := encodeSorted(w, item, indent+1); err != nil {
				return err
			}
			if i < len(typed)-1 {
				if _, err := io.WriteString(w, ","); err != nil {
					return err
				}
			}
			if _, err := io.WriteString(w, "\n"); err != nil {
				return err
			}
		}
		if len(typed) > 0 {
			writeIndent(w, indent)
		}
		_, err := io.WriteString(w, "]")
		return err
	case string:
		encoded, _ := json.Marshal(typed)
		_, err := w.Write(encoded)
		return err
	case float64, bool, nil:
		encoded, _ := json.Marshal(typed)
		_, err := w.Write(encoded)
		return err
	default:
		encoded, err := json.Marshal(typed)
		if err != nil {
			return err
		}
		_, err = w.Write(encoded)
		return err
	}
}

func writeIndent(w io.Writer, indent int) {
	for i := 0; i < indent; i++ {
		_, _ = io.WriteString(w, "  ")
	}
}
