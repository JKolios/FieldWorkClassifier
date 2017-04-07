package main

import (
	"github.com/JKolios/FieldWorkClassifier/Common/api"
	"github.com/JKolios/FieldWorkClassifier/Common/config"
	"github.com/JKolios/FieldWorkClassifier/Common/rabbitmq"
	"github.com/JKolios/FieldWorkClassifier/Common/utils"
	"github.com/streadway/amqp"
	"golang.org/x/net/context"
	"gopkg.in/olivere/elastic.v5"
	"log"
	"time"
)

func initESClient(url, username, password string, indices []string, doSniff bool) *elastic.Client {

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

func initAMQPClient(config *config.Config) (*amqp.Connection, *amqp.Channel) {

	log.Printf("Connecting to RabbitMQ on: %v", config.AmqpURL)
	conn, err := amqp.Dial(config.AmqpURL)
	utils.CheckFatalError(err)
	ch, err := conn.Channel()
	utils.CheckFatalError(err)
	log.Println("Connected to RabbitMQ.")
	for _, queue := range config.AmqpQueues {
		log.Printf("Declaring Queue: %v", queue)
		_, err = ch.QueueDeclare(
			queue,
			false, // durable
			false, // delete when unused
			false, // exclusive
			false, // no-wait
			nil,   // arguments
		)
		utils.CheckFatalError(err)
		log.Println("Queue Declared")
	}
	return conn, ch
}

func main() {

	log.Println("Starting QueryRunner")

	//Config fetch
	settings := config.GetConfFromJSONFile("config.json")

	//ES init
	esClient := initESClient(settings.ElasticURL, settings.ElasticUsername, settings.ElasticPassword, settings.Indices, settings.SniffCluster)
	defer esClient.Stop()

	//Rabbitmq init
	var amqpChannel *amqp.Channel

	if settings.UseAMQP {

		amqpConnection, amqpChannel := initAMQPClient(settings)
		defer amqpConnection.Close()
		defer amqpChannel.Close()

		rabbitmq.StartSubscribers(amqpChannel, esClient, settings)
	} else {
		amqpChannel = nil
	}

	apiInstance := api.SetupAPI(esClient, amqpChannel, settings)
	apiInstance.Run(settings.ApiURL)
}
