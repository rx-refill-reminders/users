#!/bin/bash -eu

if [[ -z "${GITHUB_REPOSITORY-}" ]]; then
    GITHUB_REPOSITORY=$(git remote show origin \
        | sed -nE 's|.*github.com:([^/]+/.*).git|\1|p' \
        | sort -u)
fi

# Parse all tags, and group them by prefix
TAGS_PARSED=$(./scripts/parse-all-tags.sh)

# Parse component info files, and identify expected prefixes for each component
COMPONENTS_PARSED=$(./scripts/component-tag-prefixes.sh)


# Match parsed tag prefixes with components
TAGS_WITH_COMPONENTS=$(jq -nrc \
    --argjson tags "$TAGS_PARSED" \
    --argjson comps "$COMPONENTS_PARSED" \
    '
        $tags
        | map(
            . as $group 
            | $group.prefix as $group_prefix
            | $comps
            | map(select(.prefix == $group_prefix))
            | . as $matching_comp
            | if ($matching_comp | length > 0) then ($group | .component = $matching_comp[0].name) else $group end
        )
        | map(select(.component != null))
    ')

echo "# Changelog"
echo ""

ITERABLE_TAGGED_COMPONENTS=$(echo "$TAGS_WITH_COMPONENTS" | jq -rc '.[]')
while IFS= read -r GROUP || [[ -n "$GROUP" ]]; do
    COMPONENT=$(echo "$GROUP" | jq -rc '.component')
    TAG_PREFIX=$(echo "$GROUP" | jq -rc '.prefix')

    NUM_VERSIONS=$(echo "$GROUP" | jq -rc '.versions | length')
    ITERABLE_VERSIONS=$(echo "$GROUP" | jq -rc '.versions | .[]' | sort -rV)
    LATEST_VERSION=$(echo "$ITERABLE_VERSIONS" | head -n 1)

    echo "## $COMPONENT"
    echo ""

    while IFS= read -r VERSION || [[ -n "$VERSION" ]]; do
        TAG="${TAG_PREFIX}${VERSION}"
        COMMIT_SHA=$(git rev-parse "tags/$TAG")
        COMMIT_SUBJECT=$(git log -1 --format=%s "$COMMIT_SHA")

        echo "* **$VERSION** ([view release](https://github.com/${GITHUB_REPOSITORY}/releases/tag/$TAG))"
        echo "<br>"
        echo "$COMMIT_SUBJECT"
        echo ""
    done <<< "$ITERABLE_VERSIONS"

    echo ""
done <<< "$ITERABLE_TAGGED_COMPONENTS"
