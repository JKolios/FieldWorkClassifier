package main

import (
	"github.com/JKolios/FieldWorkClassifier/Common/config"
	"github.com/JKolios/FieldWorkClassifier/Common/esclient"
	"github.com/JKolios/FieldWorkClassifier/QueryRunner/api"
	"log"
)

func main() {

	log.Println("Starting QueryRunner")

	//Config fetch
	settings := config.GetConfFromJSONFile("config.json")

	//ES init
	esClient := esclient.InitESClient(settings.ElasticURL,
		settings.ElasticUsername,
		settings.ElasticPassword,
		settings.SniffCluster)

	defer esClient.Stop()

	apiInstance := api.SetupAPI(esClient, settings)
	apiInstance.Run(settings.ApiURL)
}
