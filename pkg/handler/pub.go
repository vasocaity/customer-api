package handler

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/kafka-go"
)

type KafkaHandler struct {
	writer *kafka.Writer
}

func NewKafkaHandler(brokers []string, topic string) *KafkaHandler {
	return &KafkaHandler{
		writer: kafka.NewWriter(kafka.WriterConfig{
			Brokers: brokers,
			Topic:   topic,
		}),
	}
}

func (h *KafkaHandler) Close() error {
	return h.writer.Close()
}

func (h *KafkaHandler) Publish(c *gin.Context) {
	var json struct {
		Key   string `json:"key"`
		Value string `json:"value" binding:"required"`
	}

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	msg := kafka.Message{
		Key:   []byte(json.Key),
		Value: []byte(json.Value),
	}

	err := h.writer.WriteMessages(context.Background(), msg)
	if err != nil {
		log.Printf("failed to write message: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "message published"})
}
