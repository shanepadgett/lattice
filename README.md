# lcss

## Temporary CLI usage

Build and run the placeholder CLI (outputs to `./bin/lcss`):

- Build (Make, preferred): `make build`
- Build (direct fallback): `go build -o ./bin/lcss ./cmd/lcss`
- Run: `./bin/lcss`
- Version: `./bin/lcss --version`

## Releases

To create a release, tag and push a version tag that matches `v*` (for example `v1.0.0`):

- Create tag: `git tag v1.0.0`
- Push tags: `git push --tags`

Pushing a matching tag will trigger the release workflow which builds platform binaries and attaches them to the GitHub Release.

Download artifacts:

- From the GitHub Release page: open the release and click the attached binary for your platform.
- Using GitHub CLI: `gh release download <tag> --pattern "lcss-*"` (downloads release assets matching `lcss-*`).
