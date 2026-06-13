#!/bin/bash -eu

COMPONENT_INFOS=$(find .components -name 'info.yml' -mindepth 2)

COMPONENT_PREFIXES='[]'
while IFS= read -r INFO_FILE || [[ -n "$INFO_FILE" ]]; do
    # By default, the expected tag prefix is <component-name>/v
    COMPONENT_NAME=$(cat "$INFO_FILE" | yq '.component')
    EXPECTED_TAG_PREFIX="${COMPONENT_NAME}/v"

    # If a tag-prefix field is set, honor that over the default expected prefix
    COMPONENT_TAG_PREFIX=$(cat "$INFO_FILE" | yq '.["tag-prefix"]')
    if [[ "$COMPONENT_TAG_PREFIX" != "null" ]]; then
        EXPECTED_TAG_PREFIX="${COMPONENT_TAG_PREFIX}"
    fi

    COMPONENT_PREFIXES=$(echo "$COMPONENT_PREFIXES" \
        | jq -rc \
            --arg name "$COMPONENT_NAME" \
            --arg prefix "$EXPECTED_TAG_PREFIX" \
            '. += [{ name: $name, prefix: $prefix }]')
done <<< "$COMPONENT_INFOS"

printf '%s' "$COMPONENT_PREFIXES"
