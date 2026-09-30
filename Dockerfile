FROM golang:1-alpine AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/urlbit .

FROM alpine:3

RUN apk add --no-cache ca-certificates curl \
    && adduser -D -u 10001 appuser

WORKDIR /app
COPY --from=builder /out/urlbit /app/urlbit

USER appuser
EXPOSE 8082

HEALTHCHECK --interval=30s --timeout=5s --start-period=20s --retries=3 \
    CMD curl -fsS http://127.0.0.1:8082/health-check || exit 1

CMD ["/app/urlbit"]