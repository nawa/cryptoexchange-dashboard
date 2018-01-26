FROM alpine:latest
WORKDIR /app
COPY ./bin/cryptoexchange-wallet-info /app/
EXPOSE 8080

RUN apk update \
    && apk upgrade \
    && apk --no-cache add --update -t ca-certificates  \
    && apk del --purge deps; rm -rf /tmp/* /var/cache/apk/*

CMD /app/cryptoexchange-wallet-info $CRWI_ARGS