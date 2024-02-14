#!/bin/sh
set -e

rm -rf completions
mkdir completions

for sh in bash fish powershell zsh; do
  go run . completion "$sh" >"completions/dottie.$sh"
done
