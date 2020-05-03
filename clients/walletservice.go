package clients

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

const URL = "http://localhost:8081/v1"

type WalletResponse struct {
	ID        uint
	PID       string
	CreatedAt time.Time `json:"created_at"`
	User      uint      `json:"user_id"`
	Balance   uint
}

func GetWallet(id uint) (*WalletResponse, error) {
	url := fmt.Sprintf("%s/wallets/%d", URL, id)

	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("wallet not found")
	}

	target := new(WalletResponse)
	if err = json.NewDecoder(response.Body).Decode(&target); err != nil {
		return nil, err
	}
	return target, nil
}
