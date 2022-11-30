package consumer

import (
	"exam/post-service/config"
	"exam/post-service/kafka/post"
	"exam/post-service/pkg/logger"

	"github.com/jmoiron/sqlx"
)

type KafkaConsumer struct {
	ProducerCreate *post.ProducerCreateConsumer
}

type KafkaConsumerI interface {
	Post() *post.ProducerCreateConsumer
}

func NewPostconsumer(cfg config.Config, log logger.Logger, db *sqlx.DB) (KafkaConsumerI, func(), error) {
	postVer, err := post.NewProducerCreateconsumer(cfg, db)
	if err != nil {
		return &KafkaConsumer{}, func() {}, err
	}
	return &KafkaConsumer{
			ProducerCreate: postVer,
		}, func() {
			postVer.ConnClose()
		}, nil
}

func (k *KafkaConsumer) Post() *post.ProducerCreateConsumer {
	return k.ProducerCreate
}
