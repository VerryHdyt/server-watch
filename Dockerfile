FROM golang:1.24-alpine AS builder
WORKDIR /src
COPY go.mod ./
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/server-watch .

FROM alpine:3.21
RUN adduser -D -H -u 10001 appuser
COPY --from=builder /out/server-watch /usr/local/bin/server-watch
USER appuser
EXPOSE 8080
ENTRYPOINT ["server-watch"]
