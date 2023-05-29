FROM golang:1.20 as builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 go build

FROM ubuntu:latest

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/proxy .

VOLUME [ "/config/config.toml" ]

CMD [ "/app/proxy", "-config", "/config/config.toml" ]
