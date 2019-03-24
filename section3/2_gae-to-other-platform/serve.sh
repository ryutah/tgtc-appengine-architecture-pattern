#!/bin/sh

set -eux

go run ./other/main.go &

trap 'kill %1' 2 9

dev_appserver.py ./appengine
