package requests

type OfferRequest struct {
	User   uint   `json:"user_id" binding:"required"`
	Wallet uint   `json:"wallet_id" binding:"required"`
	Amount uint   `json:"amount" binding:"required"`
	From   string `json:"from" binding:"required"`
	To     string `json:"to" binding:"required"`
}
