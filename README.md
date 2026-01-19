<p align="center">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="docs/_media/lattice-logo.png">
    <source media="(prefers-color-scheme: light)" srcset="docs/_media/lattice-logo.png">
    <img alt="lattice" src="docs/_media/lattice-logo.png" style="max-width: 100%; border-radius: 6px;">
  </picture>
</p>

<p align="center">
  Lattice CSS compiler (lcss) builds utility CSS from a JSON configuration and can scan your project to understand which classes are used.
</p>

## Getting started

### Build the binary (outputs to ./bin/lcss)

```bash
make build
go build -o ./bin/lcss ./cmd/lcss
```

### Create your configuration files

Base config is required. Use configs/default.json as a starting point.
Site config is optional and lets you override the base for a specific project.

### Compile CSS

```bash
./bin/lcss build --base configs/default.json --out dist/lattice.css
```

### Include the generated CSS in your site

Add dist/lattice.css to your HTML or build pipeline.

## Configuration

- Base config: a required JSON file with tokens, utilities, and other compile rules.
- Site config: optional JSON overrides for a specific site.
- Validate and view the merged config: ./bin/lcss config print --base configs/default.json --site path/to/site.json

## Common commands

- Build CSS: ./bin/lcss build --base <path> [--site <path>] [--out <path>] [--stdout]
- Watch and rebuild: ./bin/lcss watch --base <path> [--site <path>] [--out <path>] [--interval <dur>] [--once]
- Emit tokens CSS: ./bin/lcss tokens --base <path> [--site <path>] [--out <path>] [--stdout]
- Scan project usage: ./bin/lcss scan --base <path> [--site <path>] [--top <n>]
- Version info: ./bin/lcss --version
