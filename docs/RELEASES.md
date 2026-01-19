# Releases

To create a release, tag and push a version tag that matches v* (for example v1.0.0):

- Create tag: git tag v1.0.0
- Push tags: git push --tags

Pushing a matching tag will trigger the release workflow which builds platform binaries and attaches them to the GitHub Release.

Download artifacts:

- From the GitHub Release page: open the release and click the attached binary for your platform.
- Using GitHub CLI: gh release download <tag> --pattern "lcss-*" (downloads release assets matching lcss-*).
