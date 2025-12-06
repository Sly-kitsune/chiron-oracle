FROM mcr.microsoft.com/devcontainers/go:1.22

RUN apt-get update && apt-get install -y build-essential gcc make \
 && rm -rf /var/lib/apt/lists/*

WORKDIR /workspaces/chiron-oracle

# Build Swiss Ephemeris
RUN cd swisseph && make clean && make CFLAGS="-O2 -fPIC" libswe.so \
 && cc -shared -o libswe.so swedate.o swehouse.o swejpl.o swemmoon.o swemplan.o sweph.o swephlib.o swecl.o swehel.o -lm \
 && cp libswe.so /usr/local/lib/ \
 && cp swephexp.h /usr/local/include/ \
 && ldconfig

ENV CGO_CFLAGS="-I/usr/local/include"
ENV CGO_LDFLAGS="-L/usr/local/lib -lswe -lm"
