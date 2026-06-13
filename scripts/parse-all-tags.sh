#!/bin/bash -eu

ALL_TAGS=$(git tag --list)

# Parse all valid version tags, separating the numerical semver from any characters preceding it
# v0.0.0 -> { prefix: "v", version:"0.0.0" }
VERSIONS_JSON=$(echo "$ALL_TAGS" \
    | jq -Rrsc '
        # Split input lines into an array
        split("\n")
        
        # Strip empty newlines
        | map(select(. != ""))

        # Ignore tags which are not valid versions
        | map(capture("(?<prefix>.*?)(?<version>(?:[0-9]+\\.){0,2}[0-9]+$)"))
    ')

# Organize all parsed tags by prefix
VERSIONS_BY_PREFIX=$(echo "$VERSIONS_JSON" \
    | jq -rc '
        # Group all versions by prefix
        group_by(.prefix)

        # Convert array groupings into more digestible objects (prefix, versions)
        | map({
            prefix: .[0].prefix,
            versions: (map(.version) | reverse),
        })
    ')

printf '%s' "$VERSIONS_BY_PREFIX"
