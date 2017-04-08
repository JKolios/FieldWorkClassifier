FROM golang:1.7.5
EXPOSE 8090

RUN mkdir -p /go/src/github.com/JKolios/FieldWorkClassifier/Indexer
WORKDIR /go/src/github.com/JKolios/FieldWorkClassifier/Indexer
COPY ./ .
RUN go install -v

CMD ["Indexer"]