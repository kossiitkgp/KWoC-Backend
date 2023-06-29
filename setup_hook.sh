#!/bin/sh

SOURCE_FILE="pre-commit.sh"
TARGET_FILE="${GIT_DIR:-".git"}/hooks/pre-commit"

cp "$(dirname "$0")/$SOURCE_FILE" "$TARGET_FILE"
echo "Hook set up at $TARGET_FILE"
chmod +x .git/hooks/pre-commit