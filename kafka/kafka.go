package kafka

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	kafka "github.com/segmentio/kafka-go"
)

// to produce messages
const (
	topic     = "my-topic"
	topic2    = "message-log"
	partition = 0
)

func Produce(ctx context.Context) {

	log.Println("Inside Producer...")
	conn, err := kafka.DialLeader(ctx, "tcp", "localhost:9092", topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	_, err = conn.WriteMessages(
		kafka.Message{Value: []byte("one!")},
		kafka.Message{Value: []byte("two!")},
		kafka.Message{Value: []byte("three!")},
	)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	if err := conn.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}

}

func Consume(ctx context.Context) {

	conn, err := kafka.DialLeader(ctx, "tcp", "localhost:9092", topic, partition)
	if err != nil {
		log.Fatal("failed to dial leader:", err)
	}

	//SetReadDeadline sets the deadline for future Read calls and any currently-blocked Read call. A zero value for t means Read will not time out.
	conn.SetReadDeadline(time.Now().Add(10 * time.Second))

	//ReadBatch reads a batch of messages from the kafka server.
	//The method always returns a non-nil Batch value.
	//If an error occurred, either sending the fetch request or reading the response, the error will be made available by the returned value of the batch's Close method.
	//While it is safe to call ReadBatch concurrently from multiple goroutines it may be hard for the program to predict the results as the connection offset will be read and written by multiple goroutines, they could read duplicates, or messages may be seen by only some of the goroutines.
	//A program doesn't specify the number of messages in wants from a batch, but gives the minimum and maximum number of bytes that it wants to receive from the kafka server.

	batch := conn.ReadBatch(10e3, 1e6) // fetch 10KB min, 1MB max

	b := make([]byte, 10e3) // 10KB max per message
	for {
		n, err := batch.Read(b)
		if err != nil {
			break
		}
		fmt.Println(string(b[:n]))
	}

	if err := batch.Close(); err != nil {
		log.Fatal("failed to close batch:", err)
	}

	if err := conn.Close(); err != nil {
		log.Fatal("failed to close connection:", err)
	}

}

func Produce2(ctx context.Context) {
	// initialize a counter
	i := 0

	// intialize the writer with the broker addresses, and the topic
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   topic2,
	})

	for {
		// each kafka message has a key and value. The key is used
		// to decide which partition (and consequently, which broker)
		// the message gets published on
		err := w.WriteMessages(ctx, kafka.Message{
			Key: []byte(strconv.Itoa(i)),
			// create an arbitrary message payload for the value
			Value: []byte("this is message" + strconv.Itoa(i)),
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

func Consume2(ctx context.Context) {
	// initialize a new reader with the brokers and topic
	// the groupID identifies the consumer and prevents
	// it from receiving duplicate messages
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   topic2,
		GroupID: "my-group",
	})
	for {
		// the `ReadMessage` method blocks until we receive the next event
		msg, err := r.ReadMessage(ctx)
		if err != nil {
			panic("could not read message " + err.Error())
		}
		// after receiving the message, log its value
		fmt.Println("received Topic: ", msg.Topic)
		fmt.Println("received Partition: ", msg.Partition)
		fmt.Println("received Offset: ", msg.Offset)
		fmt.Println("received Key: ", string(msg.Key))
		fmt.Println("received Value: ", string(msg.Value))
	}
}
