FROM golang

COPY cmd/vegamcache-operator/main .

WORKDIR .

CMD ["./main"]