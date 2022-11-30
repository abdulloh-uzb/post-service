package post

import (
	"context"
	"encoding/json"
	"exam/post-service/config"
	pbc "exam/post-service/genproto/customer"
	pbp "exam/post-service/genproto/post"
	"exam/post-service/pkg/logger"
	"exam/post-service/storage"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/segmentio/kafka-go"
)

type ProducerCreateConsumer struct {
	Reader    *kafka.Reader
	ConnClose func()
	Cfg       config.Config
	Logger    logger.Logger
	Storage   storage.IStorage
}

func NewProducerCreateconsumer(cfg config.Config, db *sqlx.DB) (*ProducerCreateConsumer, error) {
	r := kafka.NewReader(
		kafka.ReaderConfig{
			Brokers:   []string{"kafka:29092"},
			Topic:     cfg.ConsumerTopic,
			Partition: 0,
			MinBytes:  10e3,
			MaxBytes:  10e6,
		},
	)
	return &ProducerCreateConsumer{
		Reader:  r,
		Storage: storage.NewStorage(db),
		ConnClose: func() {
			r.Close()
		},
		Cfg: cfg,
	}, nil
}

func (p *ProducerCreateConsumer) ProducerCreateConsumerCode() {
	for {
		msg, err := p.Reader.ReadMessage(context.Background())
		if err != nil {
			fmt.Println(err)
			return
		}
		req := &pbc.Customer{}
		err = json.Unmarshal(msg.Value, req)
		if err != nil {
			fmt.Println("error while unmarsheling in kafka message", err)
			return
		}
		for _, post := range req.Posts {
			i := &pbp.PostReq{
				Id:         post.Id,
				CustomerId: post.CustomerId,
			}
			p.Storage.Post().Create(i)

		}

		fmt.Println("Message is recieved: ", string(msg.Value))
	}
}
