FROM golang:latest as builder
LABEL MAINTAINER="NanoScape Engineering <"

WORKDIR /go/src/mahir-trade-be
COPY . .

RUN go mod download && \
    go mod verify

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o /go/bin/mahir-trade-be ./cmd/mahir-trade-be

FROM debian:bullseye-slim

RUN apt-get update && \
    apt-get install -y \
    curl \
    python3 \
    python3-crcmod \
    gnupg \
    && curl -O https://dl.google.com/dl/cloudsdk/channels/rapid/downloads/google-cloud-sdk-367.0.0-linux-x86_64.tar.gz \
    && tar -xf google-cloud-sdk-367.0.0-linux-x86_64.tar.gz \
    && ./google-cloud-sdk/install.sh \
    && ./google-cloud-sdk/bin/gcloud components install gsutil \
    && rm google-cloud-sdk-367.0.0-linux-x86_64.tar.gz \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

RUN adduser --disabled-login --gecos '' appuser

RUN mkdir -p /app/temp && chown appuser:appuser /app/temp

COPY --from=builder /go/bin/mahir-trade-be /app/mahir-trade-be
COPY --from=builder /go/src/mahir-trade-be/.env /app/.env

USER appuser

WORKDIR /app

EXPOSE 8080

ENV PATH="/google-cloud-sdk/bin:$PATH"

COPY service_account.json /app/service_account.json
RUN gcloud auth activate-service-account --key-file=/app/service_account.json

ENTRYPOINT ["/app/mahir-trade-be", "-env", "/app/.env"]