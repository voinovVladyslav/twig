#!/usr/bin/env sh

set -e

go mod download

go build -o twig main.go

install -m 755 twig /usr/bin/twig
