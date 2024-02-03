#!/bin/sh
set -e

rm -rf completions
mkdir completions

for sh in bash fish ps zsh; do
  go run ./cmd/ completion "$sh" >"completions/dottie.$sh"
done
