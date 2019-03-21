#!/bin/sh

set -eux

gcloud deployment-manager deployments create \
  --config deployment.yaml \
  --automatic-rollback-on-error \
  `gcloud config get-value project`-tgtc-fileupload-example
