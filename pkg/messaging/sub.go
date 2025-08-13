package messaging

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
)

func NewSub() {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"kafka:9092"},
		Topic:     "my-topic",
		Partition: 0,
		GroupID:   "my-group",
	})
	defer reader.Close()

	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Fatal("could not read message:", err)
		}
		log.Printf("Message received: key=%s value=%s\n", string(m.Key), string(m.Value))
	}
}
