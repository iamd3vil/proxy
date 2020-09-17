FROM golang:1.15-alpine as builder

WORKDIR /app

COPY . .

RUN go build

FROM alpine:3.12

WORKDIR /app

COPY --from=builder /app/proxy .

VOLUME [ "/config/config.toml" ]

CMD [ "/app/proxy", "-config", "/config/config.toml" ]