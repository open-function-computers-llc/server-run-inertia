#!/bin/bash

killall server-run

# Don't run npx mix if SKIP_MIX is set
if [ "$SKIP_MIX" != "true" ]; then
    npm exec mix
fi


CGO_ENABLED=0 go build -ldflags="-X main.Version=`git rev-parse HEAD`" -o cmd/server-run

cmd/server-run $@ &
