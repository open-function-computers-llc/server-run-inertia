#!/bin/bash

killall server-run
npx mix

CGO_ENABLED=0 go build -o cmd/server-run

cmd/server-run $@ &
