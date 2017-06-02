FROM       alpine:3.4
MAINTAINER David Cuadrado <dacuad@facebook.com>
EXPOSE     9001

RUN apk add --update ca-certificates
COPY release/mongodb_exporter-linux-amd64 /usr/local/bin/mongodb_exporter

ENTRYPOINT [ "mongodb_exporter" ]
