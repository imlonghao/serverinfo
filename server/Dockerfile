FROM golang:1.14-alpine AS builder
LABEL maintainer="imlonghao <dockerfile@esd.cc>"
WORKDIR /builder
COPY . /builder
RUN apk add upx && \
    GO111MODULE=on go build -ldflags="-s -w" -o /app && \
    upx --lzma --best /app

FROM alpine
COPY --from=builder /app /
COPY index.html /
CMD ["/app"]