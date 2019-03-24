#!/bin/sh

set -eux

gcloud deployment-manager deployments create \
  --config deployment.yaml \
  --automatic-rollback-on-error \
  tgtc-example

if ! gcloud beta tasks queues describe tgtc-sample-queue >/dev/null 2>&1; then
  gcloud --quet beta tasks queues create-app-engine-queue tgtc-sample-queue
fi
