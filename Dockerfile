FROM golang:latest as builder
LABEL MAINTAINER="NanoScape Engineering <"

WORKDIR /go/src/mahir-trade-be
COPY . .

RUN go mod download && \
    go mod verify


RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -a -installsuffix cgo -o /go/bin/mahir-trade-be ./cmd/mahir-trade-be


FROM alpine:latest
RUN apk update && \
    adduser -D appuser

COPY --from=builder /go/bin/mahir-trade-be /app/mahir-trade-be
COPY --from=builder /go/src/mahir-trade-be/.env /app/.env

USER appuser

WORKDIR /app

EXPOSE 8080

ENTRYPOINT ["/app/mahir-trade-be", "-env", "/app/.env"]
