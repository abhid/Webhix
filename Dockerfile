FROM node:22-alpine AS ui
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM golang:1.25-alpine AS builder
WORKDIR /app
RUN apk add --no-cache gcc musl-dev
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=ui /app/internal/web/static ./internal/web/static
RUN CGO_ENABLED=1 go build -trimpath -ldflags="-s -w" -o /webhix ./cmd/webhix

FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /webhix /usr/local/bin/webhix
VOLUME ["/data"]
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=5 \
  CMD wget -q -O /dev/null http://localhost:8080/healthz || exit 1
ENTRYPOINT ["webhix", "serve"]
