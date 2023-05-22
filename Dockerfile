FROM golang:alpine as builder

RUN apk --no-cache add curl git make perl unzip
RUN wget "https://github.com/Masterminds/glide/releases/download/v0.13.3/glide-v0.13.3-linux-amd64.zip" && \
    unzip glide-v0.13.3-linux-amd64.zip && \
    mv linux-amd64/glide /usr/bin
COPY . /go/src/github.com/dcu/mongodb_exporter
RUN go env -w GO111MODULE=off
RUN cd /go/src/github.com/dcu/mongodb_exporter && make release

FROM       alpine:3.4
EXPOSE     9001

RUN apk add --update ca-certificates
COPY --from=builder /go/src/github.com/dcu/mongodb_exporter/release/mongodb_exporter-linux-amd64 /usr/local/bin/mongodb_exporter

ENTRYPOINT [ "mongodb_exporter" ]
