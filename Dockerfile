FROM golang:alpine AS builder

WORKDIR /build
ADD .env /build/
COPY . .
RUN go mod download
RUN go build -o main cmd/main.go

FROM alpine
WORKDIR /

COPY --from=builder /build/main /
COPY --from=builder /build/.env /

ENTRYPOINT ["./main"]

