#!/usr/bin/env bash
set -euo pipefail

tag=${1:-}
if [[ ! "$tag" =~ ^v(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)$ ]]; then
  echo "release version error: tag '$tag' is not strict SemVer" >&2
  exit 1
fi

version=$(sed -n 's/^const Version = "\([^"]*\)"/\1/p' version.go)
if [[ -z "$version" ]]; then
  echo "release version error: cannot parse Version from version.go" >&2
  exit 1
fi
if [[ "$tag" != "v$version" ]]; then
  echo "release version error: tag $tag != version.go v$version" >&2
  exit 1
fi

module_path=$(sed -n 's/^module //p' go.mod)
major=${version%%.*}
if (( major >= 2 )); then
  if [[ "$module_path" != */v"$major" ]]; then
    echo "release version error: Go v$major requires module path suffix /v$major" >&2
    exit 1
  fi
elif [[ "$module_path" =~ /v([2-9]|[1-9][0-9]+)$ ]]; then
  echo "release version error: Go v$major must not use module path ${BASH_REMATCH[0]}" >&2
  exit 1
fi

echo "release version $version is consistent with module $module_path"
