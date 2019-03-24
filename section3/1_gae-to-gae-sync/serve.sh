#!/bin/sh

set -eux

dev_appserver.py \
  --env_var GOOGLE_CLOUD_PROJECT=${GOOGLE_CLOUD_PROJECT} \
  --env_var DEV_SERVER=true \
  ./front-service ./micro-service ./dispatch.yaml
