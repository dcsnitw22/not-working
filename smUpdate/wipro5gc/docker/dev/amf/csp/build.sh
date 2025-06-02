#!/bin/bash

#cd "$(dirname "$0")"

set -euo pipefail

buildArch=`uname -m`
case "${buildArch##*-}" in
          aarch64) ;;
        x86_64) ;;
        *) echo "Current architecture (${buildArch}) is not supported."; exit 2; ;;
esac

echo "==============================================="
#echo " Image: ${AMF_CSP_IMAGE_TAG:=w5gc_amf_csp}"
echo " Image: ${IMAGE_TAG:="$1"}"
echo "==============================================="

set -x

sudo docker build -f $(dirname "$0")/Dockerfile \
          --tag ${IMAGE_TAG} \
 ${DOCKER_BUILD_ARGS-} .

docker run --rm "${IMAGE_TAG}" csp -h || true

