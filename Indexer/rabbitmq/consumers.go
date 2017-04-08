package rabbitmq

import (
	"github.com/JKolios/FieldWorkClassifier/Indexer/es"
	"github.com/JKolios/FieldWorkClassifier/Common/config"
	"github.com/JKolios/FieldWorkClassifier/Common/utils"
	"github.com/streadway/amqp"
	"gopkg.in/olivere/elastic.v5"
	"log"
)

func StartSubscribers(amqpChan *amqp.Channel, esClient *elastic.Client, config *config.Settings) {
	log.Println("Starting RabbitMQ subscribers")
	msgChan, err := amqpChan.Consume(config.AmqpQueues[0], "", true, false, false, false, nil)
	utils.CheckFatalError(err)
	go incomingDocConsumer(msgChan, esClient, config)
	log.Println("Started RabbitMQ subscribers")
}

func incomingDocConsumer(incomingChan <-chan amqp.Delivery, esClient *elastic.Client, config *config.Settings) {
	for message := range incomingChan {
		log.Printf("Received incoming Doc: %s", message.Body)
		resp, err := es.IndexDocJSONBytes(esClient, config.DefaultIndex, "document", string(message.Body))
		log.Println(resp)
		utils.CheckFatalError(err)
		log.Println("Indexed Incoming Doc")
	}
}
