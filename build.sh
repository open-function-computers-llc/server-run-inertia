#!/bin/bash

killall server-run

# Don't run npx mix if app is running in development mode
if [ "$APP_ENV" != "development" ]; then
    npm exec mix
fi


CGO_ENABLED=0 go build -ldflags="-X main.Version=`git rev-parse HEAD`" -o cmd/server-run

cmd/server-run $@ &
