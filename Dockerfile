FROM golang:1.22

RUN apt-get update && apt-get install -y build-essential gcc make \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY . .

# Build Swiss Ephemeris manually
RUN cd swisseph \
    && cc -O2 -fPIC -c *.c \
    && cc -shared -o libswe.so *.o -lm \
    && cp libswe.so /usr/local/lib/ \
    && cp swephexp.h /usr/local/include/ \
    && ldconfig

ENV CGO_CFLAGS="-I/usr/local/include"
ENV CGO_LDFLAGS="-L/usr/local/lib -lswe -lm"

RUN go mod download
RUN go build -ldflags="-w -s" -o out

CMD ["./out"]
