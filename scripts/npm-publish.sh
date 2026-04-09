#!/bin/bash

set -e

if [ "$CI" != "true" ]; then
    echo "Not running in CI, skipping commit"
    exit 0
fi

echo "Publishing to npm registry"
npm publish --access public || {
    echo "Failed to publish to npm registry"
    exit 1
}
