#FROM golang:latest as builder
FROM karalabe/xgo-1.7.x as builder

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
WORKDIR /go/src/github.com/mpppk/unravel-twitter-worker
COPY . .
RUN make

# runtime image
FROM alpine
RUN apk add --no-cache ca-certificates
COPY --from=builder /go/src/github.com/mpppk/unravel-twitter-worker /unravel-twitter-worker
ENTRYPOINT ["/unravel-twitter-worker"]