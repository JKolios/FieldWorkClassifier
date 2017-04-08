FROM golang:1.7.5
EXPOSE 8091

WORKDIR /go/src/app
COPY . .
RUN ls -al

RUN go install -v

CMD ["app"]