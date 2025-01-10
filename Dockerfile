FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o emqx-s3-backup

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/emqx-s3-backup .

ENV EMQX_URL=""
ENV EMQX_API_NAME=""
ENV EMQX_API_PASS=""
ENV S3_BUCKET=""
ENV S3_REGION=""
ENV S3_PATH=""
USER nobody

CMD ["./emqx-s3-backup"]
