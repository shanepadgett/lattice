package schema

import (
	"encoding/json"
	"errors"
	"fmt"

	"lcss/internal/version"

	"os"
	"path/filepath"
)

func Generate(versionString string, schemaVersion int) error {
	if schemaVersion <= 0 {
		return errors.New("--schema-version is required and must be greater than zero")
	}

	if major, ok := version.ParseSemverMajor(versionString); ok {
		if major != schemaVersion {
			return fmt.Errorf("--schema-version must match --version major (%d)", major)
		}
	}

	latestID := "https://raw.githubusercontent.com/shanepadgett/lattice/main/schemas/lattice.schema.json"
	releaseID := ""
	if major, ok := version.ParseSemverMajor(versionString); ok && major > 0 {
		releaseID = fmt.Sprintf("https://github.com/shanepadgett/lattice/releases/download/%s/lattice.schema.json", versionString)
	}

	schema := buildRootSchema(latestID, releaseID, schemaVersion)

	latestBytes, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return err
	}
	latestBytes = append(latestBytes, '\n')

	if err := os.MkdirAll("schemas", 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join("schemas", "lattice.schema.json"), latestBytes, 0o644); err != nil {
		return err
	}

	if err := os.MkdirAll("dist", 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join("dist", "lattice.schema.json"), latestBytes, 0o644); err != nil {
		return err
	}

	return nil
}

type schemaDoc map[string]any

func buildRootSchema(latestID, releaseID string, schemaVersion int) schemaDoc {
	ids := []any{latestID}
	if releaseID != "" {
		ids = append(ids, releaseID)
	}

	return schemaDoc{
		"$schema":              "https://json-schema.org/draft/2020-12/schema",
		"$id":                  ids,
		"title":                "Lattice configuration",
		"type":                 "object",
		"additionalProperties": false,
		"required":             []any{"schemaVersion", "themes", "scales"},
		"properties": schemaDoc{
			"schemaVersion": schemaDoc{
				"type":  "integer",
				"const": schemaVersion,
			},
			"classPrefix": schemaDoc{"type": "string"},
			"separator":   schemaDoc{"type": "string"},
			"breakpoints": schemaDoc{
				"type":                 "object",
				"additionalProperties": schemaDoc{"type": "string"},
			},
			"themes": schemaDoc{
				"type":                 "object",
				"additionalProperties": themeSchema(),
				"required":             []any{"default"},
			},
			"fonts": fontsSchema(),
			"scales": schemaDoc{
				"type":                 "object",
				"additionalProperties": false,
				"required":             []any{"space"},
				"properties":           scalesSchema(),
			},
			"variants": variantsSchema(),
			"build":    buildBuildSchema(),
		},
	}
}

func themeSchema() schemaDoc {
	return schemaDoc{
		"type":                 "object",
		"additionalProperties": false,
		"properties": schemaDoc{
			"colors": schemaDoc{
				"type":                 "object",
				"additionalProperties": schemaDoc{"type": "string"},
			},
			"font": schemaDoc{
				"type":                 "object",
				"additionalProperties": schemaDoc{"type": "string"},
			},
		},
	}
}

func fontsSchema() schemaDoc {
	return schemaDoc{
		"type":                 "object",
		"additionalProperties": false,
		"properties": schemaDoc{
			"imports": schemaDoc{
				"type":  "array",
				"items": schemaDoc{"type": "string"},
			},
			"faces": schemaDoc{
				"type":  "array",
				"items": fontFaceSchema(),
			},
		},
	}
}

func fontFaceSchema() schemaDoc {
	return schemaDoc{
		"type":                 "object",
		"additionalProperties": false,
		"required":             []any{"family", "src"},
		"properties": schemaDoc{
			"family":            schemaDoc{"type": "string"},
			"style":             schemaDoc{"type": "string"},
			"weight":            schemaDoc{"type": "string"},
			"stretch":           schemaDoc{"type": "string"},
			"display":           schemaDoc{"type": "string"},
			"unicodeRange":      schemaDoc{"type": "string"},
			"featureSettings":   schemaDoc{"type": "string"},
			"variationSettings": schemaDoc{"type": "string"},
			"src": schemaDoc{
				"type":     "array",
				"minItems": 1,
				"items": schemaDoc{
					"type":                 "object",
					"additionalProperties": false,
					"required":             []any{"url"},
					"properties": schemaDoc{
						"url":    schemaDoc{"type": "string"},
						"format": schemaDoc{"type": "string"},
						"tech":   schemaDoc{"type": "string"},
					},
				},
			},
		},
	}
}

func scalesSchema() schemaDoc {
	props := schemaDoc{}
	props["space"] = scaleMapSchema()
	props["size"] = scaleMapSchema()
	props["radius"] = scaleMapSchema()
	props["borderWidth"] = scaleMapSchema()
	props["fontSize"] = scaleMapSchema()
	props["lineHeight"] = scaleMapSchema()
	props["fontWeight"] = scaleMapSchema()
	props["letterSpacing"] = scaleMapSchema()
	props["shadow"] = scaleMapSchema()
	props["z"] = scaleMapSchema()
	props["opacity"] = scaleMapSchema()
	props["aspect"] = scaleMapSchema()
	props["duration"] = scaleMapSchema()
	props["easing"] = scaleMapSchema()
	props["delay"] = scaleMapSchema()
	props["translate"] = scaleMapSchema()
	props["rotate"] = scaleMapSchema()
	props["scale"] = scaleMapSchema()
	props["maxWidth"] = scaleMapSchema()
	props["maxHeight"] = scaleMapSchema()
	props["container"] = scaleMapSchema()
	return props
}

func scaleMapSchema() schemaDoc {
	return schemaDoc{
		"type":                 "object",
		"additionalProperties": schemaDoc{"type": "string"},
	}
}

func variantsSchema() schemaDoc {
	return schemaDoc{
		"type":                 "object",
		"additionalProperties": false,
		"properties": schemaDoc{
			"responsive": schemaDoc{
				"type":  "array",
				"items": schemaDoc{"type": "string"},
			},
			"state": schemaDoc{
				"type":  "array",
				"items": schemaDoc{"type": "string"},
			},
		},
	}
}

func buildBuildSchema() schemaDoc {
	return schemaDoc{
		"type":                 "object",
		"additionalProperties": false,
		"properties": schemaDoc{
			"content":  schemaDoc{"type": "array", "items": schemaDoc{"type": "string"}},
			"safelist": schemaDoc{"type": "array", "items": schemaDoc{"type": "string"}},
			"gridColumns": schemaDoc{
				"type":    "integer",
				"minimum": 0,
			},
			"unknownClassPolicy": schemaDoc{
				"type": "string",
				"enum": []any{"ignore", "warn", "error"},
			},
			"emit": schemaDoc{
				"type":                 "object",
				"additionalProperties": false,
				"properties": schemaDoc{
					"fontsCss":  schemaDoc{"type": "boolean"},
					"tokensCss": schemaDoc{"type": "boolean"},
					"base":      schemaDoc{"type": "boolean"},
					"manifest":  schemaDoc{"type": "boolean"},
				},
			},
		},
	}
}
