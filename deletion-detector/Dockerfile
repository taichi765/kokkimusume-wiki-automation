FROM golang:latest AS builder

WORKDIR /build/
COPY ./deletion-detector/ ./deletion-detector/
COPY ./wikiwiki/ ./wikiwiki/
RUN cd deletion-detector&&CGO_ENABLED=0 go build

FROM alpine:latest

COPY --from=builder /build/deletion-detector/deletion-detector /opt/deletion-detector/deletion-detector

ENTRYPOINT [ "/opt/deletion-detector/deletion-detector" ]