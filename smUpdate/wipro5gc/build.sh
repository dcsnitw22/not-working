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
#echo " Image: ${SMF_PDUSMSP_IMAGE_TAG:=w5gc_smf_pdusmsp}"
echo "==============================================="

set -x

docker build -f ${DOCKERFILE_DIR}/Dockerfile \
          --tag ${IMAGE_TAG} \
 ${DOCKER_BUILD_ARGS-} .

docker run --rm "${IMAGE_TAG}" ${BINARY_NAME} -h || true
