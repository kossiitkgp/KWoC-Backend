#!/bin/sh

SOURCE_FILE="pre-commit.sh"
TARGET_FILE=".git/hooks/pre-commit"

cp "$(dirname "$0")/$SOURCE_FILE" "$(dirname "$0")/$TARGET_FILE"