FROM golang:latest AS builder
RUN apt-get update && apt-get install -y xz-utils && rm -rf /var/lib/apt/lists/*
ADD https://github.com/upx/upx/releases/download/v3.96/upx-3.96-amd64_linux.tar.xz /usr/local
RUN xz -d -c /usr/local/upx-3.96-amd64_linux.tar.xz | tar -xOf - upx-3.96-amd64_linux/upx > /bin/upx && chmod a+x /bin/upx
RUN mkdir /go/src/wb-l0
WORKDIR /go/src/wb-l0
COPY . .
RUN go mod download
WORKDIR /go/src/wb-l0/cmd/subscriber
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/src/wb-l0/bin/subscriber .
WORKDIR /go/src/wb-l0
RUN strip --strip-unneeded ./bin/subscriber
RUN upx ./bin/subscriber

WORKDIR /go/src/wb-l0/cmd/publisher
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/src/wb-l0/bin/publisher .
WORKDIR /go/src/wb-l0
RUN strip --strip-unneeded ./bin/publisher
RUN upx ./bin/publisher

FROM alpine:latest AS subscriber
WORKDIR /root
COPY --from=builder /go/src/wb-l0/bin/subscriber .
CMD ["./subscriber"]

FROM alpine:latest AS publisher
WORKDIR /root
COPY --from=builder /go/src/wb-l0/bin/publisher .
CMD ["./publisher"]

