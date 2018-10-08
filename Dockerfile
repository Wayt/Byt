FROM golang:latest

MAINTAINER max@wayt.me

RUN apt-get update -qq && \
    apt-get install -qqy npm ruby-sass
RUN ln -s /usr/bin/nodejs /usr/bin/node

WORKDIR /go/src/github.com/byttl/byt/
ADD . .

RUN npm install && \
    npm install -g grunt-cli
RUN grunt

RUN go build

FROM alpine:latest
COPY --from=0 /go/src/github.com/wayt/byt/byt /byt

RUN mkdir -p /data
VOLUME /data
ENV BYT_UPLOAD_DIR /data
ENV BYT_HOST "localhost:8080"
ENV BYT_BIND ":8080"

EXPOSE 8080

CMD ["/byt"]

