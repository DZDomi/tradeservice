package prices

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"strings"
	"time"
)

const schema = "wss"
const address = "dex.binance.org"
const path = "/api/ws/$all@allTickers"

type SubscribeMessage struct {
	Method  string   `json:"method"`
	Topic   string   `json:"topic"`
	Symbols []string `json:"symbols"`
}

type TickerMessage struct {
	Data []struct {
		Symbol      string `json:"s"`
		PriceChange string `json:"p"`
		Bid         string `json:"b"`
		Ask         string `json:"a"`
	} `json:"data"`
}

func Subscribe() {
	u := url.URL{Scheme: schema, Host: address, Path: path}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			var response *TickerMessage
			err = json.Unmarshal(message, &response)
			if err != nil {
				fmt.Println(err)
			}
			for _, tick := range response.Data {
				if strings.Index(tick.Symbol, "BTC") >= 0 && strings.Index(tick.Symbol, "USD") >= 0 {
					log.Printf("%s: %s, %s/%s", tick.Symbol, tick.PriceChange, tick.Ask, tick.Bid)
				}
			}
		}
	}()

	//err = c.WriteJSON(&SubscribeMessage{
	//	Method:  "subscribe",
	//	Topic:   "marketDepth",
	//	Symbols: []string{"BNB_ETH"},
	//})
	//if err != nil {
	//	fmt.Println("Binance error: " + err.Error())
	//}

	for {
		time.Sleep(time.Second * 5)
	}

}
