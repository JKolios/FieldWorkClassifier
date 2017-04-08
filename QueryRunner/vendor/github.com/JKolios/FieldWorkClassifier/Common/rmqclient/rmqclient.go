package rmqclient

import (
	"github.com/streadway/amqp"
	"log"
	"github.com/JKolios/FieldWorkClassifier/Common/config"
	"github.com/JKolios/FieldWorkClassifier/Common/utils"
)

func InitAMQPClient(config *config.Settings) (*amqp.Connection, *amqp.Channel) {


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