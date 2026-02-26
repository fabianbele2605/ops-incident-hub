FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copiar go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copiar código fuente
COPY . .

# Compilar aplicación
RUN go build -o bin/api ./cmd/api

# Imagen final
FROM alpine:latest

WORKDIR /app

# Instalar ca-certificates para HTTPS
RUN apk --no-cache add ca-certificates

# Copiar binario y migraciones
COPY --from=builder /app/bin/api .
COPY --from=builder /app/migrations ./migrations

# Exponer puerto
EXPOSE 8080

# Ejecutar aplicación
CMD ["./api"]
