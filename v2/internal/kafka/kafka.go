package kafka

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
	"log"
	"strconv"
	"time"
)

const ScanningTopic = "start-scanning"
const Group = "scanning"

const kafkaBroker = ""
const kafkaUsername = ""
const kafkaPassword = ""

//var kafkaBroker = os.Getenv("KAFKA_BROKER")
//var kafkaUsername = os.Getenv("KAFKA_USERNAME")
//var kafkaPassword = os.Getenv("KAFKA_PASSWORD")

func Produce(ctx context.Context, kafkaDialer *kafka.Dialer) {
	// initialize a counter
	i := 0

	// initialize the writer with the broker addresses, and the topic
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{kafkaBroker},
		Topic:   ScanningTopic,
		Dialer:  kafkaDialer,
	})

	for {
		// each kafka message has a key and value. The key is used
		// to decide which partition (and consequently, which broker)
		// the message gets published on
		err := w.WriteMessages(ctx, kafka.Message{
			Key: []byte(strconv.Itoa(i)),
			// create an arbitrary message payload for the value
			Value: []byte("this is message " + strconv.Itoa(i)),
		})
		if err != nil {
			panic("could not write message " + err.Error())
		}

		// log a confirmation once the message is written
		fmt.Println("writes:", i)
		i++
		// sleep for a second
		time.Sleep(time.Second)
	}
}

func GetReader() *kafka.Reader {
	// initialize a new reader with the brokers and topic
	// the groupID identifies the consumer and prevents
	// it from receiving duplicate messages
	dialer := GetDialer()
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkaBroker},
		GroupID: Group,
		Topic:   ScanningTopic,
		Dialer:  dialer,
	})
	//for {
	//	// the `ReadMessage` method blocks until we receive the next event
	//	msg, err := r.ReadMessage(ctx)
	//	if err != nil {
	//		panic("could not read message " + err.Error())
	//	}
	//	// after receiving the message, log its value
	//	fmt.Println("received: ", string(msg.Key), string(msg.Value))
	//	scanner.StartScanning("ngocrongonline.com", "risk-id-vip")
	//}
	return r
}

func GetDialer() *kafka.Dialer {
	// get kafka reader using environment variables.
	mechanism, err := scram.Mechanism(
		scram.SHA256,
		kafkaUsername,
		kafkaPassword,
	)

	if err != nil {
		log.Fatalln(err)
	}
	dialer := &kafka.Dialer{
		SASLMechanism: mechanism,
		TLS:           &tls.Config{},
	}

	return dialer
}
