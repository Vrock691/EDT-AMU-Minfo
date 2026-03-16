FROM golang:1.24.6-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY src ./src

RUN CGO_ENABLED=0 GOOS=linux go build -o /out/server ./src

FROM alpine:3.20

RUN apk add --no-cache ca-certificates

WORKDIR /app/src

COPY --from=builder /out/server /app/server
COPY src/html ./html

EXPOSE 10000

CMD ["/app/server"]
