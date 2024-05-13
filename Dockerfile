# Gunakan image Golang resmi sebagai base image
FROM golang:latest

# Set environment variable agar Go menggunakan mode production
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Buat direktori kerja di dalam container
WORKDIR /app

# Salin file go.mod dan go.sum terlebih dahulu dan lakukan download dependensi
COPY go.mod .
COPY go.sum .
RUN go mod download

# Salin seluruh kode sumber aplikasi
COPY . .
COPY .env .

# Build aplikasi Golang
RUN go build -o mahir-trade-be ./cmd/mahir-trade-be

# Expose port yang digunakan oleh aplikasi
EXPOSE 8080

# Atur command untuk dijalankan saat container dijalankan
CMD ["go", "run", "cmd/mahir-trade-be/main.go"]
