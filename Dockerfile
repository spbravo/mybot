FROM golang:1.15

WORKDIR /go/src/mybot
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

ENTRYPOINT ["mybot"]
