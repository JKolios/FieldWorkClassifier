package main

import (
	"github.com/JKolios/FieldWorkClassifier/Common/config"
	"github.com/JKolios/FieldWorkClassifier/Common/esclient"
	"github.com/JKolios/FieldWorkClassifier/Indexer/api"
	"log"
	"github.com/JKolios/FieldWorkClassifier/Indexer/es"

)


func main() {

	log.Println("Starting Indexer")

	//Config fetch
	log.Println("Loading configuration")
	settings := config.GetConfFromJSONFile("config.json")

	//ES client init
	elasticClient := esclient.InitESClient(settings.ElasticURL,
		settings.ElasticUsername,
		settings.ElasticPassword,
		settings.SniffCluster)
	defer elasticClient.Stop()

	//Create the required indices and set their mappings
	es.InitIndices(elasticClient)

	//Initialize the field coordinate storage doc
	es.InitFieldLocationDocument(elasticClient)

	//Add percolator queries to the device_data index
	es.InitPercolators(elasticClient)

	//Create the HTTP and WS endpoints and listen for connections
	apiInstance := api.SetupAPI(elasticClient, settings)
	apiInstance.Run(settings.ApiURL)
}
