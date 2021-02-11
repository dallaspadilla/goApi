FROM golang

RUN mkdir -p /go/src/goApi

WORKDIR /go/src/goApi

COPY . /go/src/goApi

RUN go install goApi

CMD /go/bin/goApi

EXPOSE 8080