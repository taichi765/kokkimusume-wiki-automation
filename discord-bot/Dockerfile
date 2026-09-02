FROM golang:latest AS builder

ARG version
ARG commitHash

WORKDIR /build/
COPY ./discord-bot/ ./discord-bot/
COPY ./types/ ./types/
RUN cd discord-bot&&CGO_ENABLED=0 go build -ldflags \
    "-X main.version=${version} -X main.commitHash=${commitHash}"

FROM alpine:latest

COPY --from=builder /build/discord-bot/discord-bot /opt/discord-bot/discord-bot

EXPOSE 8080

ENTRYPOINT [ "/opt/discord-bot/discord-bot" ]
