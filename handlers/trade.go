package handlers

import (
	"fmt"
	"github.com/DZDomi/tradeservice/clients"
	"github.com/DZDomi/tradeservice/models"
	"github.com/DZDomi/tradeservice/requests"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/google/uuid"
	"net/http"
	"time"
)

func CreateOffer(c *gin.Context) {
	var request requests.OfferRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	from := &models.Asset{}
	to := &models.Asset{}
	models.DB.Where(&models.Asset{
		Name: request.From,
	}).First(from)
	models.DB.Where(&models.Asset{
		Name: request.To,
	}).First(to)

	if from.Name == "" || to.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Assets do not exist"})
		return
	}

	fromWalletResponse, err := clients.GetWallet(request.FromWallet)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	toWalletResponse, err := clients.GetWallet(request.ToWallet)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if fromWalletResponse.User != request.User || toWalletResponse.User != request.User {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid wallet"})
		return
	}

	if fromWalletResponse.Balance < request.Amount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient balance"})
		return
	}

	pid, err := uuid.NewRandom()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("Generated Offer: " + pid.String())

	// TODO Calculation logic

	offer := &models.Trade{
		PID:        pid,
		CreatedAt:  time.Now(),
		From:       from.Name,
		To:         to.Name,
		User:       request.User,
		FromWallet: request.FromWallet,
		ToWallet:   request.ToWallet,
	}

	if err = clients.SetObject("offer", pid.String(), offer, time.Minute); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, offer)
}

func Accept(c *gin.Context) {
	pid := c.Param("id")

	// TODO: Think about how this can work with long execution times
	lock, err := clients.GetLock(pid, time.Minute)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer clients.ReleaseLock(lock)

	offer := &models.Trade{}
	if err := clients.GetObject("offer", pid, offer); err == redis.Nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No offer found for this id"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	now := time.Now()
	offer.Accepted = &now

	models.DB.Create(offer)
	if offer.ID == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create trade"})
		return
	}

	if err := clients.DeleteObject("offer", pid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("Sending offer:", offer.PID, "to kafka...")
	if err = clients.TradeCreated(offer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, offer)
}

func ListTrades(c *gin.Context) {
	userId := c.Query("user_id")

	trades := new([]models.Trade)
	if userId != "" {
		models.DB.Where("user_id = ?", userId).Find(&trades)
	} else {
		models.DB.Find(&trades)
	}

	c.JSON(http.StatusOK, trades)
}
