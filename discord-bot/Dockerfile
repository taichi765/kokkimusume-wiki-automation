FROM golang:latest AS builder

WORKDIR /build/
COPY ./discord-bot/ ./discord-bot/
COPY ./common ./common/
RUN cd discord-bot&&CGO_ENABLED=0 go build

FROM alpine:latest

COPY --from=builder /build/discord-bot/discord-bot /opt/discord-bot/discord-bot

EXPOSE 8080

ENTRYPOINT [ "/opt/discord-bot/discord-bot" ]
