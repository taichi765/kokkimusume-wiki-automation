FROM golang:latest AS builder

WORKDIR /build/
COPY ./discord-bot/ ./discord-bot/
COPY ./common ./common/
RUN cd discord-bot&&go build -o /build/discord-bot

FROM alpine:latest

COPY --from=builder /build/discord-bot /opt/discord-bot/discord-bot

ENTRYPOINT [ "/opt/discord-bot/discord-bot" ]
