package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/DZDomi/tradeservice/models"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"log"
	"time"
)

var topicsToListen = [...]string{
	"wallets-updated",
}

var topicsToWrite = [...]string{
	"trade-created",
}

var listeners = map[string]*kafka.Reader{}
var writers = map[string]*kafka.Writer{}

type TradeCreatedMessage struct {
	ID         uint   `json:"id"`
	User       uint   `json:"user_id"`
	FromWallet uint   `json:"from_wallet_id"`
	ToWallet   uint   `json:"to_wallet_id"`
	From       string `json:"from"`
	To         string `json:"to"`
	Amount     uint   `json:"amount"`
}

type Wallet struct {
	User   uint   `json:"user_id"`
	Wallet uint   `json:"wallet_id"`
	Amount uint   `json:"amount"`
	Action string `json:"action"`
}

type TriggerEvent struct {
	ID   uint   `json:"id"`
	Type string `json:"type"`
}

type WalletsUpdatedEvent struct {
	TriggeredBy *TriggerEvent `json:"triggered_by"`
	Wallets     []Wallet      `json:"wallets"`
}

func InitKafka() {
	for _, topic := range topicsToListen {
		listeners[topic] = kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{"localhost:9092"},
			GroupID: "trade",
			Topic:   topic,
		})
		go listenForTopic(topic)
	}
	for _, topic := range topicsToWrite {
		writers[topic] = kafka.NewWriter(kafka.WriterConfig{
			Brokers:  []string{"localhost:9092"},
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		})
	}
}

func listenForTopic(topic string) {
	for {
		fmt.Println("Listening for topic", topic)
		m, err := listeners[topic].ReadMessage(context.Background())
		if err != nil {
			log.Fatalln(err)
		}
		switch topic {
		case "wallets-updated":
			go settleTrade(m)
		}
		fmt.Printf("message at topic:%v partition:%v offset:%v	%s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
	}
}

func settleTrade(message kafka.Message) {
	walletsUpdatedEvent := &WalletsUpdatedEvent{}
	if err := json.Unmarshal(message.Value, walletsUpdatedEvent); err != nil {
		fmt.Println("Unable to decode message:", message.Value)
		return
	}

	if walletsUpdatedEvent.TriggeredBy.Type != "trade" {
		fmt.Println("Event not for me, skipping...")
		return
	}

	//TODO: Probably needs a lock

	trade := &models.Trade{}
	models.DB.Where("id = ?", walletsUpdatedEvent.TriggeredBy.ID).First(trade)

	now := time.Now()
	trade.Executed = &now

	models.DB.Save(trade)

}

func sendToTopic(topic string, message string) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	// TODO: Think about the error handling here
	go writers[topic].WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(id.String()),
		Value: []byte(message),
	})
	//if err != nil {
	//	return err
	//}
	return nil
}

func TradeCreated(trade *models.Trade) error {
	message, err := json.Marshal(&TradeCreatedMessage{
		ID:         trade.ID,
		User:       trade.User,
		FromWallet: trade.FromWallet,
		ToWallet:   trade.ToWallet,
		From:       trade.From,
		To:         trade.To,
		//TODO: Amount
		Amount: 100,
	})
	if err != nil {
		return err
	}

	if err := sendToTopic("trade-created", string(message)); err != nil {
		return err
	}
	return nil
}
