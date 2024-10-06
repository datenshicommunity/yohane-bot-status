# Gunakan base image golang:1.23 dengan rootless
FROM golang:1.23-alpine AS builder

# Buat user non-root
RUN adduser -D -g '' appuser

# Set working directory
WORKDIR /app

# Salin file go.mod dan go.sum
COPY go.mod go.sum ./

# Download dependensi
RUN go mod download

# Salin kode sumber
COPY . .

# Build aplikasi
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Tahap akhir
FROM alpine:latest

# Salin binary dari builder
COPY --from=builder /app/main /app/main

# Salin user non-root dari builder
COPY --from=builder /etc/passwd /etc/passwd

# Ganti ke user non-root
USER appuser

# Jalankan aplikasi
CMD ["/app/main"]


