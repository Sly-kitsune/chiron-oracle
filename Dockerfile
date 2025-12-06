# Use official Go image as base
FROM golang:1.22

# Install build tools
RUN apt-get update && apt-get install -y build-essential gcc make \
    && rm -rf /var/lib/apt/lists/*

# Set working directory
WORKDIR /app

# Copy repo into container
COPY . .

# Build Swiss Ephemeris
RUN cd swisseph && make clean && make CFLAGS="-O2 -fPIC" libswe.so \
    && cc -shared -o libswe.so swedate.o swehouse.o swejpl.o swemmoon.o swemplan.o sweph.o swephlib.o swecl.o swehel.o -lm \
    && cp libswe.so /usr/local/lib/ \
    && cp swephexp.h /usr/local/include/ \
    && ldconfig

# Set CGO flags so Go can link against libswe
ENV CGO_CFLAGS="-I/usr/local/include"
ENV CGO_LDFLAGS="-L/usr/local/lib -lswe -lm"

# Download Go modules
RUN go mod download

# Build your app
RUN go build -ldflags="-w -s" -o out

# Run the app
CMD ["./out"]
