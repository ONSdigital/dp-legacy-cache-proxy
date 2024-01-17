#!/bin/bash -eux

pushd dp-legacy-cache-proxy
  make build
  cp build/dp-legacy-cache-proxy Dockerfile.concourse ../build
popd
