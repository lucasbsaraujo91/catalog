package kafkahelper

import "github.com/segmentio/kafka-go"

func GetKafkaWriter(brokerAddress, topic string) *kafka.Writer {
	return kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{brokerAddress},
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})
}
