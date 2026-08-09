FROM golang:latest AS builder

RUN go build -o /opt/kokkimusume_wiki .

FROM alpine:latest

COPY --from=builder /opt/kokkimusume_wiki /opt/kokkimusume_wiki

ENTRYPOINT [ "/opt/kokkimusume_wiki" ]

