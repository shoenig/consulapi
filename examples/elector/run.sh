#!/usr/bin/env bash

set -e -u -o pipefail

go clean
go build

./elector 1