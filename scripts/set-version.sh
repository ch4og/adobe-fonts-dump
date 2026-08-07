#!/usr/bin/env bash
set -euo pipefail

version=${1:?usage: $0 VERSION}

if [[ ! $version =~ ^[0-9]+(\.[0-9]+){2}([.-][0-9A-Za-z.-]+)?$ ]]; then
    printf 'invalid version: %s\n' "$version" >&2
    exit 1
fi

root=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)

update() {
    local file=$1
    local pattern=$2
    local replacement=$3
    local backup="$file.bak"

    sed -E -i.bak "s|$pattern|$replacement|" "$file"
    rm -f "$backup"
}

update "$root/src/main.go" \
    '^const version = "[^"]+"' \
    "const version = \"$version\""
update "$root/Makefile" \
    '^VERSION \?= .*' \
    "VERSION ?= $version"

printf 'Updated version to %s\n' "$version"
