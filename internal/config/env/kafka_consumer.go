package env

import (
	"errors"
	"os"
	"strings"

	"github.com/IBM/sarama"
)

const (
	brokersEnvName = "KAFKA_BROKERS"
	groupIDEnvName = "KAFKA_GROUP_ID"
)

type kafkaConsumerConfig struct {
	brokers []string
	groupID string
}

// NewKafkaConsumerConfig Конструктор конфига слушателя.
func NewKafkaConsumerConfig() (*kafkaConsumerConfig, error) {
	brokersStr := os.Getenv(brokersEnvName)
	if len(brokersStr) == 0 {
		return nil, errors.New("kafka brokers address not found")
	}

	brokers := strings.Split(brokersStr, ",")

	groupID := os.Getenv(groupIDEnvName)
	if len(groupID) == 0 {
		return nil, errors.New("kafka group id not found")
	}

	return &kafkaConsumerConfig{
		brokers: brokers,
		groupID: groupID,
	}, nil
}

func (c *kafkaConsumerConfig) Brokers() []string {
	return c.brokers
}

func (c *kafkaConsumerConfig) GroupID() string {
	return c.groupID
}

// Config возвращает конфигурацию для sarama consumer
func (c *kafkaConsumerConfig) Config() *sarama.Config {
	config := sarama.NewConfig()
	config.Version = sarama.V2_6_0_0
	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	return config
}
