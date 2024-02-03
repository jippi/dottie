#!/bin/sh
set -e

rm -rf manpages
mkdir manpages

go run ./cmd | gzip -c -9 >manpages/dotti.1.gz
