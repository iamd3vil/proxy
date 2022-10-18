FROM golang:1.19-alpine as builder

RUN apk add git --no-cache

WORKDIR /app

COPY . .

RUN go build

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/proxy .

VOLUME [ "/config/config.toml" ]

CMD [ "/app/proxy", "-config", "/config/config.toml" ]