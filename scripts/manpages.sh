#!/bin/sh
set -e

rm -rf manpages
mkdir manpages

go run ./cmd | gzip -c -9 >manpages/dottie.1.gz
