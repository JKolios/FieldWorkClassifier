FROM golang:1.7.5
EXPOSE 8091

RUN mkdir -p /go/src/github.com/JKolios/FieldWorkClassifier/QueryRunner
WORKDIR /go/src/github.com/JKolios/FieldWorkClassifier/QueryRunner
COPY ./ .
RUN go install -v

CMD ["QueryRunner"]