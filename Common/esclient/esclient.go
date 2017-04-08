package esclient

import (
	"log"
	"time"
	"github.com/JKolios/FieldWorkClassifier/Common/utils"
	"gopkg.in/olivere/elastic.v5"

)

func InitESClient(url, username, password string, doSniff bool) *elastic.Client {

	log.Printf("Connecting to ES on: %v", url)
	retrier := elastic.NewBackoffRetrier(elastic.NewSimpleBackoff(1000))

	elasticClient, err := elastic.NewClient(
		elastic.SetURL(url),
		elastic.SetSniff(doSniff),
		elastic.SetRetrier(retrier),
		elastic.SetSnifferTimeout(time.Second*30),
		elastic.SetBasicAuth(username, password))

	utils.CheckFatalError(err)

	log.Println("Connected to ES")

	return elasticClient
}