#!/bin/sh

set -eux

dev_appserver.py api front dispatch.yaml
