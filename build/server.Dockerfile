FROM golang:alpine as build

WORKDIR /opt

COPY go.mod go.sum ./
RUN  go mod download

COPY cmd/server cmd/server
COPY internal   internal

RUN cd /opt/cmd/server && \
    go build -o /srv/server


FROM alpine:latest

COPY --from=build /srv /srv

WORKDIR /srv
CMD /srv/server
