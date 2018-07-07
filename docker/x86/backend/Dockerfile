FROM golang:1.10 as builder

WORKDIR /go/src/github.com/nawa/cryptoexchange-dashboard
COPY . /go/src/github.com/nawa/cryptoexchange-dashboard

RUN CGO_ENABLED=0 GOOS=linux make build

FROM alpine:latest
WORKDIR /app
COPY --from=builder /go/src/github.com/nawa/cryptoexchange-dashboard/bin/* /app/
EXPOSE 8080

RUN apk update \
    && apk upgrade \
    && apk --no-cache add --update -t ca-certificates  \
    && apk del --purge deps; rm -rf /tmp/* /var/cache/apk/*

CMD /app/cryptoexchange-dashboard $CREXD_ARGS