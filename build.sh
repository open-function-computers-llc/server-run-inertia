#!/bin/bash

killall server-run
npx mix

CGO_ENABLED=0 go build -o httpd/server-run

httpd/server-run &
