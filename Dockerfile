FROM golang:1.17-alpine AS builder

RUN apk add make
WORKDIR /build
COPY . .
RUN go get
RUN make build

FROM alpine:latest

WORKDIR /app/nino/timeouts
COPY --from=builder /build/timeouts /app/nino/timeouts/timeouts
CMD ["/app/nino/timeouts/timeouts"]
