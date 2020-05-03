package clients

import (
	"context"
	"encoding/json"
	"github.com/DZDomi/tradeservice/models"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

var tradeCreatedWriter *kafka.Writer

type TradeCreatedMessage struct {
	ID     uuid.UUID `json:"id"`
	User   uint      `json:"user_id"`
	Wallet uint      `json:"wallet_id"`
	From   string    `json:"from"`
	To     string    `json:"to"`
	Amount uint      `json:"amount"`
}

func InitKafka() {
	tradeCreatedWriter = kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    "trade-created",
		Balancer: &kafka.LeastBytes{},
	})
}

func TradeCreated(trade *models.Trade) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	message, err := json.Marshal(&TradeCreatedMessage{
		ID:     id,
		User:   trade.User,
		Wallet: trade.Wallet,
		From:   trade.From,
		To:     trade.To,
		//TODO: Amount
		//amount:   ,
	})
	if err != nil {
		return err
	}
	return tradeCreatedWriter.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(id.String()),
		Value: message,
	})
}
