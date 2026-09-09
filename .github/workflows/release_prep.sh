#!/usr/bin/env bash
# Builds the release archive and prints the release notes to stdout.
#
# Called by bazel-contrib/.github/.github/workflows/release_ruleset.yaml, which
# requires this exact path and reads the release notes from stdout.
set -o errexit -o nounset -o pipefail

TAG="$1"
VERSION="${TAG#v}"
REPO_NAME="upspin-gdrive"

# The registry entry points at this archive. `.bcr/source.template.json` gives
# `strip_prefix` as `{REPO}-{VERSION}`, so it holds one top level directory
# named after the version, which is what GitHub's own generated tarball did
# before this script replaced it.
#
# The entry used to point at `archive/refs/tags/{TAG}.tar.gz`, the tarball
# GitHub generates on request. That one cannot be attested, because no
# workflow produces it, and BCR flags it as an unstable url because GitHub
# does not promise its bytes stay the same. This builds a release asset
# instead.
git archive --format=tar.gz --prefix="${REPO_NAME}-${VERSION}/" \
  -o "${TAG}.tar.gz" HEAD

cat <<NOTES
## Using Bzlmod

\`\`\`starlark
bazel_dep(name = "upspin_drive", version = "${VERSION}")
\`\`\`
NOTES
