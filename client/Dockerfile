FROM golang:1.14-alpine AS builder
LABEL maintainer="imlonghao <dockerfile@esd.cc>"
WORKDIR /builder
COPY . /builder
RUN apk add upx && \
    GO111MODULE=on go build -ldflags="-s -w" -o /app && \
    upx --lzma --best /app

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app /
CMD ["/app"]