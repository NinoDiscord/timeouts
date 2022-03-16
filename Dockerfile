FROM golang:1.18-alpine AS builder

RUN apk add make jq git
WORKDIR /
COPY . .
RUN go get
RUN make build

FROM alpine:latest

WORKDIR /app/nino/timeouts
COPY --from=builder /build/timeouts /app/nino/timeouts/timeouts
CMD ["/app/nino/timeouts/timeouts"]
