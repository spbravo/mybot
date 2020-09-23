FROM golang:1.15.2

WORKDIR /mybot/src/mybot
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

ENTRYPOINT ["mybot"]
