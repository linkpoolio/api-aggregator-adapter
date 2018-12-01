FROM golang:1.10-alpine as builder

RUN apk add --no-cache make curl git gcc musl-dev linux-headers
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

ADD . /go/src/github.com/linkpoolio/api-aggregator-cl-ea
RUN cd /go/src/github.com/linkpoolio/api-aggregator-cl-ea && go get && go build -o api-aggregator-cl-ea main.go

# Copy adaptor into a second stage container
FROM alpine:latest

RUN apk add --no-cache ca-certificates
COPY --from=builder /go/src/github.com/linkpoolio/api-aggregator-cl-ea/api-aggregator-cl-ea /usr/local/bin/

EXPOSE 8080
ENTRYPOINT ["api-aggregator-cl-ea"]