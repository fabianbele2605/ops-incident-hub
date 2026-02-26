# Build stage
FROM golang:1.22-alpine AS builder

# Instalar dependencias de build
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /build

# Copiar y descargar dependencias (cache layer)
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copiar código fuente
COPY . .

# Compilar con optimizaciones
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o api ./cmd/api

# Runtime stage
FROM scratch

# Metadata
LABEL maintainer="ops-incident-hub" \
      version="1.0" \
      description="Ops Incident Hub API"

# Copiar archivos necesarios desde builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

WORKDIR /app

# Copiar binario y migraciones
COPY --from=builder /build/api .
COPY --from=builder /build/migrations ./migrations

# Exponer puerto
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD ["/app/api", "health"] || exit 1

# Ejecutar aplicación
ENTRYPOINT ["/app/api"]
CMD []
