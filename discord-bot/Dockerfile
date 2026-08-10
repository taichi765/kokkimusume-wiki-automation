FROM golang:latest AS builder

RUN go build -o /opt/discord-bot/

FROM alpine:latest

COPY --from=builder /opt/discord-bot/discord-bot /opt/discord-bot/discord-bot

ENTRYPOINT [ "/opt/discord-bot/discord-bot" ]
