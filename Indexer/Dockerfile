FROM golang:1.7.5
EXPOSE 8090

WORKDIR /go/src/app
COPY . .
RUN ls -al

RUN go install -v

CMD ["app"]