package requests

type OfferRequest struct {
	User       uint   `json:"user_id" binding:"required"`
	FromWallet uint   `json:"from_wallet_id" binding:"required"`
	ToWallet   uint   `json:"to_wallet_id" binding:"required"`
	Amount     uint   `json:"amount" binding:"required"`
	From       string `json:"from" binding:"required"`
	To         string `json:"to" binding:"required"`
}
