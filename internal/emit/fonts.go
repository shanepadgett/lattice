package emit

import (
	"strings"

	"lcss/internal/config"
)

func FontsCSS(cfg config.Config) []byte {
	if len(cfg.Fonts.Imports) == 0 && len(cfg.Fonts.Faces) == 0 {
		return nil
	}

	var b strings.Builder

	for _, url := range cfg.Fonts.Imports {
		if strings.TrimSpace(url) == "" {
			continue
		}
		b.WriteString("@import url(\"")
		b.WriteString(url)
		b.WriteString("\");\n")
	}

	for _, face := range cfg.Fonts.Faces {
		if strings.TrimSpace(face.Family) == "" || len(face.Src) == 0 {
			continue
		}
		b.WriteString("@font-face {\n")
		b.WriteString("  font-family: \"")
		b.WriteString(face.Family)
		b.WriteString("\";\n")
		if strings.TrimSpace(face.Style) != "" {
			b.WriteString("  font-style: ")
			b.WriteString(face.Style)
			b.WriteString(";\n")
		}
		if strings.TrimSpace(face.Weight) != "" {
			b.WriteString("  font-weight: ")
			b.WriteString(face.Weight)
			b.WriteString(";\n")
		}
		if strings.TrimSpace(face.Stretch) != "" {
			b.WriteString("  font-stretch: ")
			b.WriteString(face.Stretch)
			b.WriteString(";\n")
		}
		if strings.TrimSpace(face.Display) != "" {
			b.WriteString("  font-display: ")
			b.WriteString(face.Display)
			b.WriteString(";\n")
		}
		b.WriteString("  src: ")
		b.WriteString(buildFontSources(face.Src))
		b.WriteString(";\n")
		if strings.TrimSpace(face.UnicodeRange) != "" {
			b.WriteString("  unicode-range: ")
			b.WriteString(face.UnicodeRange)
			b.WriteString(";\n")
		}
		if strings.TrimSpace(face.FeatureSettings) != "" {
			b.WriteString("  font-feature-settings: ")
			b.WriteString(face.FeatureSettings)
			b.WriteString(";\n")
		}
		if strings.TrimSpace(face.VariationSettings) != "" {
			b.WriteString("  font-variation-settings: ")
			b.WriteString(face.VariationSettings)
			b.WriteString(";\n")
		}
		b.WriteString("}\n")
	}

	return []byte(b.String())
}

func buildFontSources(sources []config.FontSource) string {
	parts := make([]string, 0, len(sources))
	for _, source := range sources {
		if strings.TrimSpace(source.URL) == "" {
			continue
		}
		var b strings.Builder
		b.WriteString("url(\"")
		b.WriteString(source.URL)
		b.WriteString("\")")
		if strings.TrimSpace(source.Format) != "" {
			b.WriteString(" format(\"")
			b.WriteString(source.Format)
			b.WriteString("\")")
		}
		if strings.TrimSpace(source.Tech) != "" {
			b.WriteString(" tech(\"")
			b.WriteString(source.Tech)
			b.WriteString("\")")
		}
		parts = append(parts, b.String())
	}
	return strings.Join(parts, ", ")
}
