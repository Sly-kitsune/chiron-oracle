FROM mcr.microsoft.com/devcontainers/go:1.22

RUN apt-get update && apt-get install -y build-essential gcc make \
 && rm -rf /var/lib/apt/lists/*

WORKDIR /workspaces/chiron-oracle

ENV CGO_CFLAGS="-I/usr/local/include"
ENV CGO_LDFLAGS="-L/usr/local/lib -lswe -lm"
