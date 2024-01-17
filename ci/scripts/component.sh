#!/bin/bash -eux

pushd dp-legacy-cache-proxy
  make test-component
popd
