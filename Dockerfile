FROM       alpine:3.4
MAINTAINER David Cuadrado <dacuad@facebook.com>
EXPOSE     9001

COPY . /go/src/github.com/dcu/mongodb_exporter

RUN apk add --update -t build-deps go git make curl \
    && export GOPATH=/go && export PATH=$GOPATH/bin:$PATH \
    && cd $GOPATH/src/github.com/dcu/mongodb_exporter \
    && mkdir $GOPATH/bin && curl https://glide.sh/get | sh && make deps \
    && go build -o /bin/mongodb_exporter \
    && apk del --purge build-deps && rm -rf $GOPATH

ENTRYPOINT [ "/bin/mongodb_exporter" ]
