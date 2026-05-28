#!/bin/bash
IMAGE_NAME="phoenix-ebpf-builder"
CONTAINER_DIR="/Users/fallofpheonix/os/pheonixos/phoenix_os/telemetry/ebpf"

# Build the docker image if it doesn't exist
docker build -t $IMAGE_NAME -f build.dockerfile .

# Run the compilation inside the container
docker run --rm \
    -v "$CONTAINER_DIR:/build" \
    -v "/sys/kernel/btf:/sys/kernel/btf:ro" \
    $IMAGE_NAME \
    make
