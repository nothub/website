#!/usr/bin/env bash

set -eu

mkdir -p "content"
mkdir -p "posts"

for path in "content"/*; do
  echo >&2 "processing: ${path}"

  slug="$(basename "${path%.*}" | inline-detox)"
  echo >&2 "slug: ${slug}"

  if [[ -d "${path}" ]]; then
    for asset_path in "${path}"/*; do
      mkdir -p "posts/${slug}"
      asset_file="$(basename "${asset_path}")"
      if [[ "${asset_file}" == "index.md" ]]; then
        continue
      fi
      cp "${asset_path}" "posts/${slug}/${asset_file}"
    done
    path="${path}/index.md"
  fi

  if [[ ${path} != *.md ]]; then
    echo >&2 "ignoring non-markdown: ${path}"
    continue
  fi

  # styles: breezedark espresso haddock kate monochrome pygments tango zenburn
  pandoc \
    --fail-if-warnings \
    --highlight-style monochrome \
    --number-sections \
    --table-of-contents \
    --toc-depth=2 \
    --from=markdown \
    --to=html \
    --output "posts/${slug}.html" \
    "${path}"

  echo >&2 "done"
done
