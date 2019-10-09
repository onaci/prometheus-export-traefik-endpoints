FROM golang:1-alpine as builder

RUN apk --update upgrade \
    && apk --no-cache --no-progress add git make gcc musl-dev

WORKDIR /src

ENV GO111MODULE on
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN make build

FROM alpine:3.10
RUN apk --update upgrade \
    && apk --no-cache --no-progress add ca-certificates

COPY --from=builder /src/prometheus-export-traefik-endpoints /usr/bin/

ENTRYPOINT ["/usr/bin/prometheus-export-traefik-endpoints"]
