# Gunakan image Go resmi sebagai base image
FROM golang:1.23-alpine AS builder

# Atur direktori kerja
WORKDIR /app

# Salin file go.mod dan go.sum
COPY go.mod go.sum ./

# Download dependensi
RUN go mod download

# Salin kode sumber
COPY . .

# Build aplikasi
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Gunakan image Alpine yang lebih kecil untuk image akhir
FROM alpine:latest

# Atur direktori kerja
WORKDIR /root/

# Salin binary dari builder stage
COPY --from=builder /app/main .

# Jalankan aplikasi
CMD ["./main"]
