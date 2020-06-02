FROM golang:1.13.7 AS builder
COPY . /src
RUN cd /src && go build -o goapp

FROM alpine:3.12.0
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
WORKDIR /app
COPY --from=builder /src/goapp /app/
ENTRYPOINT ./goapp