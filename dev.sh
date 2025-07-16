#!/bin/bash

killall server-run
killall "npm exec mix watch"

APP_LIVERELOAD_PORT=8983 go run wrappers/liveReload.go $@