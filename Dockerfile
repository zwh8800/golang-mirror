FROM golang:1.6.2-alpine
MAINTAINER zwh8800 <496781108@qq.com>

WORKDIR $GOPATH/src/github.com/zwh8800/golang-mirror

RUN apk update && apk add ca-certificates && apk add git && \
    apk add tzdata && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && go get github.com/Masterminds/glide && \
    go get github.com/shadowsocks/shadowsocks-go/cmd/shadowsocks-local

ADD . $GOPATH/src/github.com/zwh8800/golang-mirror/

RUN glide install && go build

VOLUME $GOPATH/src/github.com/zwh8800/golang-mirror/log
VOLUME $GOPATH/src/github.com/zwh8800/golang-mirror/config
VOLUME $GOPATH/src/github.com/zwh8800/golang-mirror/ws

CMD ["./golang-mirror", "-log_dir", "log", "-alsologtostderr", "-config", "config/golang-mirror.gcfg"]
