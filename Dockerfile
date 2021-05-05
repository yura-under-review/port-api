FROM golang:1.16.3-buster as builder

WORKDIR /app
COPY . .

RUN make build

FROM debian:buster-slim

COPY --from=builder /app/artifacts/svc /
COPY --from=builder /app/static-html/root.html /

EXPOSE 8080

WORKDIR /

CMD ["./svc"]