# 개발 스테이지
FROM golang:1.24-alpine AS dev
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
CMD ["go", "run", "cmd/main.go"]

# 빌드 스테이지
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/main.go

# 프로덕션 스테이지
FROM alpine:latest AS prod
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
CMD ["./main"]