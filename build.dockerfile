FROM --platform=linux/amd64 ubuntu:22.04

RUN apt-get update && apt-get install -y \
    clang \
    llvm \
    libbpf-dev \
    linux-headers-generic \
    libelf-dev \
    zlib1g-dev \
    make \
    pkg-config \
    gcc-multilib \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /build
