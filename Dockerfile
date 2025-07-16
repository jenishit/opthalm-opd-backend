# ---- build stage ----
FROM golang:1.26-alpine AS builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /out/server ./cmd/main.go

# ---- runtime stage ----
FROM alpine:3.20
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /out/server ./server

EXPOSE 8082
ENTRYPOINT [ "./server" ]