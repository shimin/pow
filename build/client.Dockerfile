FROM golang:alpine as build

WORKDIR /opt

COPY go.mod go.sum ./
RUN  go mod download

COPY cmd/client cmd/client
COPY internal   internal

RUN cd /opt/cmd/client && \
    go build -o /srv/client


FROM alpine:latest

COPY --from=build /srv /srv

WORKDIR /srv
CMD /srv/client
