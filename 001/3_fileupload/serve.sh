#!/bin/sh

set -eux

dev_appserver.py \
  --env_var GOOGLE_CLOUD_PROJECT=${GOOGLE_CLOUD_PROJECT} \
  --env_var GCS_BUCKET=${GCS_BUCKET} \
  .
