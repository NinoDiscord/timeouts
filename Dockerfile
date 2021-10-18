FROM golang:1.17-alpine AS builder
WORKDIR /
COPY . .
RUN go get
RUN go build

FROM alpine:latest
WORKDIR /
COPY --from=builder /timeouts /app/timeouts
CMD ["/app/timeouts"]
