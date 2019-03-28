FROM golang:1.11
COPY . /go/src/better-cgd/
WORKDIR /go/src/better-cgd/
RUN go install -ldflags="-s -w" ./...
ENTRYPOINT better-cgd