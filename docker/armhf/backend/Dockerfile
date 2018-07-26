FROM golang:1.10 as builder

WORKDIR /go/src/github.com/nawa/cryptoexchange-dashboard
COPY . /go/src/github.com/nawa/cryptoexchange-dashboard

RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 make build


FROM container4armhf/armhf-alpine:latest

WORKDIR /app
COPY docker/armhf/qemu-arm-static /usr/bin
COPY --from=builder /go/src/github.com/nawa/cryptoexchange-dashboard/bin/* /app/
EXPOSE 8080

RUN apk update \
    && apk upgrade \
    && apk --no-cache add --update -t ca-certificates  \
    && apk del --purge deps; rm -rf /tmp/* /var/cache/apk/* \
    && rm /usr/bin/qemu-arm-static

CMD /app/cryptoexchange-dashboard $CREXD_ARGS
