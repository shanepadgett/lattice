package schema

import (
	"encoding/json"
	"errors"
	"fmt"

	"lcss/internal/version"

	"os"
	"path/filepath"
)

func Generate(versionString string) error {
	releaseID := ""
	if _, ok := version.ParseSemverMajor(versionString); ok {
		releaseID = fmt.Sprintf("https://github.com/shanepadgett/lattice/releases/download/%s/lattice.schema.json", versionString)
	}
	if releaseID == "" {
		return errors.New("--version must be a valid semver tag like v1.2.3")
	}

	data, err := os.ReadFile(filepath.Join("configs", "lattice.schema.json"))
	if err != nil {
		return err
	}

	var schema map[string]any
	if err := json.Unmarshal(data, &schema); err != nil {
		return err
	}

	schema["$id"] = releaseID

	latestBytes, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return err
	}
	latestBytes = append(latestBytes, '\n')

	if err := os.MkdirAll("dist", 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join("dist", "lattice.schema.json"), latestBytes, 0o644); err != nil {
		return err
	}

	return nil
}
