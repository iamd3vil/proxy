FROM golang:1.20 as builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 go build

FROM ubuntu:latest

WORKDIR /app

COPY --from=builder /app/proxy .

VOLUME [ "/config/config.toml" ]

CMD [ "/app/proxy", "-config", "/config/config.toml" ]
