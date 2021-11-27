#!/bin/sh
go get code.rocketnine.space/tslocum/godoc-static
go get golang.org/x/tools/cmd/godoc
mkdir -p ./html
godoc-static \
    -site-name="DataQ" \
    -site-description-file=../DESCRIPTION.md \
    -destination=./html \
    ..\