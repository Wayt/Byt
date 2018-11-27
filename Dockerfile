FROM golang:latest

MAINTAINER max@wayt.me

WORKDIR /go/src/github.com/wayt/byt/
ADD . .

RUN go build -o byt

FROM golang:latest
COPY --from=0 /go/src/github.com/wayt/byt/byt /byt

RUN mkdir -p /data
VOLUME /data
ENV BYT_UPLOAD_DIR /data
ENV BYT_HOST "localhost:8080"
ENV BYT_BIND ":8080"

EXPOSE 8080

CMD ["/byt"]

