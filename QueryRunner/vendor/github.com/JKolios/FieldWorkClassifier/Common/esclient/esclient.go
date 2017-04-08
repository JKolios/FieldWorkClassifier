package esclient

import (
	"log"
	"time"
	"github.com/JKolios/FieldWorkClassifier/Common/utils"
	"gopkg.in/olivere/elastic.v5"
	"golang.org/x/net/context"

)

func InitESClient(url, username, password string, indices []string, doSniff bool) *elastic.Client {

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

	for _, index := range indices {

		log.Printf("Initializing Index: %s", index)

		indexExists, err := elasticClient.IndexExists(index).Do(context.TODO())
		utils.CheckFatalError(err)
		if !indexExists {
			resp, err := elasticClient.CreateIndex(index).Do(context.TODO())
			utils.CheckFatalError(err)
			if !resp.Acknowledged {
				log.Fatalf("Cannot create index: %s on ES", index)
			}
			log.Printf("Created index: %s on ES", index)

		} else {
			log.Printf("Index: %s already exists on ES", index)
		}

		_, err = elasticClient.OpenIndex(index).Do(context.TODO())
		utils.CheckFatalError(err)

		mapping, err := elasticClient.GetMapping().Index(index).Do(context.TODO())
		if err != nil {
			log.Printf("Cannot get mapping for index: %s", index)
		}
		log.Printf("Mapping for index %s: %s", index, mapping)
	}

	return elasticClient
}