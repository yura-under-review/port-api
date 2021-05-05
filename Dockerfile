FROM golang:1.16.3-buster as builder

WORKDIR /app
COPY . .

RUN make build

FROM debian:buster-slim

COPY --from=builder /app/artifacts/svc /

EXPOSE 8080

WORKDIR /

CMD ["./svc"]