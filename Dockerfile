FROM golang:alpine3.9 AS builder

WORKDIR /go/src/github.com/mvochoa/graphqldoc
COPY ./ ./

RUN apk add --update git \
    && go get github.com/go-bindata/go-bindata/... \
    && cd cmd/ \
    && GO111MODULE=on go install -v

FROM alpine:3.8

COPY --from=builder /go/bin/cmd /usr/local/bin/graphqldoc

CMD ["graphqldoc"]