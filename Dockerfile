FROM golang:1.24.4 AS builder

WORKDIR /app

COPY cmd /app/cmd
COPY internal /app/internal
COPY go.mod /app/go.mod
COPY go.sum /app/go.sum

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gollery ./cmd/gollery/main.go

FROM alpine:3.20

WORKDIR /gollery

COPY --from=builder /app/gollery ./gollery

COPY migrations/ /gollery/migrations/

COPY LICENSE /gollery/LICENSE

EXPOSE 8000

CMD ["./gollery"]