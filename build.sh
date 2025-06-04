#!/bin/bash

killall server-run
npx mix

CGO_ENABLED=0 go build -ldflags="-X main.Version=`git rev-parse HEAD`" -o cmd/server-run

cmd/server-run $@ &
