package main

import (
	"github.com/JKolios/FieldWorkClassifier/Common/config"
	"github.com/JKolios/FieldWorkClassifier/Common/esclient"
	"github.com/JKolios/FieldWorkClassifier/Common/rmqclient"
	"github.com/JKolios/FieldWorkClassifier/Indexer/rabbitmq"
	"github.com/JKolios/FieldWorkClassifier/Indexer/api"
	"github.com/streadway/amqp"
	"log"
)


func main() {

	log.Println("Starting Indexer")

	//Config fetch
	settings := config.GetConfFromJSONFile("config.json")

	//ES init
	esClient := esclient.InitESClient(settings.ElasticURL, settings.ElasticUsername, settings.ElasticPassword, settings.Indices, settings.SniffCluster)
	defer esClient.Stop()

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

	apiInstance := api.SetupAPI(esClient, amqpChannel, settings)
	apiInstance.Run(settings.ApiURL)
}
