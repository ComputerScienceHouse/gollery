FROM alpine:3.22 AS builder

RUN apk add npm go make

WORKDIR /app

COPY . /app

RUN CGO_ENABLED=0 GOOS=linux make build

FROM alpine:3.20

WORKDIR /gollery

COPY --from=builder /app/dist/gollery ./gollery

COPY migrations/ /gollery/migrations/

COPY LICENSE /gollery/LICENSE

EXPOSE 8000

CMD ["./gollery"]
