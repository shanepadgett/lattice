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

### Configuration defaults

Lattice ships with sensible defaults built into the binary.
If lattice.json exists in your project root, it is merged as overrides.
Use configs/default.json as a reference when authoring your own overrides.

### Compile CSS

```bash
./bin/lcss build
# OR
./bin/lcss build --out <directory>/<name>.css
```

### Include the generated CSS in your site

Add dist/lattice.css to your HTML or build pipeline.

## Configuration

- Base config: embedded defaults with tokens, utilities, and compile rules.
- Site config: optional lattice.json or --site overrides for a specific site.
- Validate and view the merged config: ./bin/lcss config print [--site path/to/site.json]

## Common commands

- Build CSS: ./bin/lcss build [--site <path>] [--out <path>] [--stdout]
- Watch and rebuild: ./bin/lcss watch [--site <path>] [--out <path>] [--interval <dur>] [--once]
- Emit tokens CSS: ./bin/lcss tokens [--site <path>] [--out <path>] [--stdout]
- Scan project usage: ./bin/lcss scan [--site <path>] [--top <n>]
- Version info: ./bin/lcss --version
