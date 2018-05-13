FROM golang

ADD . /go/src/github.com/schoolboy/vegamcacheexample

WORKDIR /go/src/github.com/schoolboy/vegamcacheexample

RUN go get ./

RUN go build main.go

EXPOSE 90 8000 8087

CMD ["./main"]