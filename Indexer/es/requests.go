package es

import (
	"code.google.com/p/go-uuid/uuid"
	"golang.org/x/net/context"
	"gopkg.in/olivere/elastic.v5"
)


func IndexDocJSONBytes(client *elastic.Client, indexName, docType string, body string) (*elastic.IndexResponse, error) {
	resp, err := client.Index().Index(indexName).Type(docType).Id(uuid.New()).BodyString(body).Do(context.TODO())
	return resp, err
}