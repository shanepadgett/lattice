package emit

import (
	"errors"
	"os"
	"path/filepath"
)

type Artifacts struct {
	LatticeCSS []byte
	Manifest   []byte
}

func Write(artifacts Artifacts, outPath string) error {
	if outPath == "" {
		return errors.New("output path is required")
	}
	if len(artifacts.LatticeCSS) == 0 && len(artifacts.Manifest) == 0 {
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return err
	}
	if len(artifacts.LatticeCSS) > 0 {
		if err := os.WriteFile(outPath, artifacts.LatticeCSS, 0o644); err != nil {
			return err
		}
	}
	if len(artifacts.Manifest) > 0 {
		manifestPath := filepath.Join(filepath.Dir(outPath), "manifest.json")
		if err := os.WriteFile(manifestPath, artifacts.Manifest, 0o644); err != nil {
			return err
		}
	}

	return nil
}
