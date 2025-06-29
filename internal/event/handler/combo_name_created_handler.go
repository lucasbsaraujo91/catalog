package handler

import (
	"catalog/pkg/events"
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/segmentio/kafka-go"
)

type ComboNameCreatedHandler struct {
	writer *kafka.Writer
}

func NewComboNameCreatedHandler(writer *kafka.Writer) *ComboNameCreatedHandler {
	return &ComboNameCreatedHandler{writer: writer}
}

func (h *ComboNameCreatedHandler) Handle(event events.EventInterface, wg *sync.WaitGroup) {
	defer wg.Done()

	jsonOutput, err := json.Marshal(event.GetPayload())
	if err != nil {
		log.Printf("Erro ao serializar payload: %v\n", err)
		return
	}

	msg := kafka.Message{
		Key:   []byte("comboname"),
		Value: jsonOutput,
	}

	err = h.writer.WriteMessages(context.Background(), msg)
	if err != nil {
		log.Printf("Erro ao enviar mensagem para o Kafka: %v\n", err)
	} else {
		log.Println("✅ Mensagem enviada com sucesso para o Kafka")
	}
}
