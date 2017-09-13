# stage 1 
FROM golang:1.8

RUN mkdir -p /go/src/github.com/dcu/mongodb_exporter && \
  curl https://glide.sh/get | sh
COPY . /go/src/github.com/dcu/mongodb_exporter/

WORKDIR /go/src/github.com/dcu/mongodb_exporter

RUN make release 

# stage 2
FROM       alpine:3.4
EXPOSE     9001

RUN apk add --update ca-certificates
COPY --from=0 /go/src/github.com/dcu/mongodb_exporter/release/mongodb_exporter-linux-amd64 /usr/local/bin/mongodb_exporter

ENTRYPOINT [ "mongodb_exporter" ]
