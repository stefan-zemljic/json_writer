#!/usr/bin/bash

set -euo pipefail

git add .
git amend
git fetch --tags
highest_tag=$(git tag | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+$' | sort -V | tail -n 1)
git tag | xargs git push -d origin
git tag | xargs git tag -d
next_tag=$(echo "$highest_tag" | awk -F. -v OFS=. '{$NF++; print}')
git tag "$next_tag"
git push origin "$next_tag"