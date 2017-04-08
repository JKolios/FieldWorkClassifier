package main

import (
	"github.com/JKolios/FieldWorkClassifier/Common/config"
	"github.com/JKolios/FieldWorkClassifier/Common/esclient"
	"github.com/JKolios/FieldWorkClassifier/Common/rmqclient"
	"github.com/JKolios/FieldWorkClassifier/Indexer/rabbitmq"
	"github.com/JKolios/FieldWorkClassifier/Indexer/api"
	"github.com/streadway/amqp"
	"log"
	"github.com/JKolios/FieldWorkClassifier/Indexer/es"
)


func main() {

	log.Println("Starting Indexer")

	//Config fetch
	settings := config.GetConfFromJSONFile("config.json")

	//ES client init
	esClient := esclient.InitESClient(settings.ElasticURL,
		settings.ElasticUsername,
		settings.ElasticPassword,
		settings.SniffCluster)
	defer esClient.Stop()

	//Create the required indices and set their mappings
	es.InitIndices(esClient)

	//Rabbitmq init
	var amqpChannel *amqp.Channel

	if settings.UseAMQP {

		amqpConnection, amqpChannel := rmqclient.InitAMQPClient(settings)
		defer amqpConnection.Close()
		defer amqpChannel.Close()

		rabbitmq.StartSubscribers(amqpChannel, esClient, settings)
	} else {
		amqpChannel = nil
	}

	//Create the HTTP and WS endpoints and listen for connections
	apiInstance := api.SetupAPI(esClient, amqpChannel, settings)
	apiInstance.Run(settings.ApiURL)
}
