package kafka

import (
	"encoding/json"
	"fmt"

	"github.com/ExchangeRates/produce-mock/internal/model"
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

var (
	TYPE_HEADER = "__TypeId__"
	CUP_TYPE    = "cup"
)

type ProducerCupRatePoint struct {
	log    *logrus.Logger
	topic  string
	urls   []string
	config *sarama.Config
}

func NewProducerCupRatePoint(urls []string) *ProducerCupRatePoint {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Retry.Max = 5
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Net.SASL.Enable = true
	config.Net.SASL.Mechanism = sarama.SASLTypePlaintext
	config.Net.SASL.User = "admin"
	config.Net.SASL.Password = "admin-secret"

	return &ProducerCupRatePoint{
		log:    logrus.New(),
		topic:  "parsed.CUP",
		urls:   urls,
		config: config,
	}
}

func (p *ProducerCupRatePoint) ProduceAll(points []model.CupRatePoint) error {
	conn, err := sarama.NewSyncProducer(p.urls, p.config)
	if err != nil {
		return err
	}
	defer conn.Close()

	for _, point := range points {
		if err := p.produce(conn, point); err != nil {
			return err
		}
	}
	p.log.WithFields(logrus.Fields{
		"size": len(points),
		"type": CUP_TYPE,
	}).Info("Events was send to broker")

	return nil
}

func (p *ProducerCupRatePoint) produce(conn sarama.SyncProducer, point model.CupRatePoint) error {
	topic := p.constructTopic(point.Major, point.Minor)
	pointBytes, err := json.Marshal(point)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(pointBytes),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte(TYPE_HEADER),
				Value: []byte(CUP_TYPE),
			},
		},
	}

	_, _, err = conn.SendMessage(msg)
	if err != nil {
		return err
	}

	return nil
}

func (p *ProducerCupRatePoint) constructTopic(major, minor string) string {
	return fmt.Sprintf("%s-%s-%s", p.topic, major, minor)
}
