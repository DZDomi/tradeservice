package handlers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/google/uuid"
	"net/http"
	"time"
	"tradeservice/clients"
	"tradeservice/models"
	"tradeservice/requests"
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

	// Check amount and wallet to user

	pid, err := uuid.NewRandom()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Println("Generated Offer: " + pid.String())

	// TODO Calculation logic

	offer := &models.Trade{
		PID:       pid,
		CreatedAt: time.Now(),
		From:      from.Name,
		To:        to.Name,
		User:      request.User,
	}

	if err = clients.SetObject("offer", pid.String(), offer, time.Minute); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"offer": offer})
}

func Execute(c *gin.Context) {
	pid := c.Param("id")

	// TODO: Think about how this can work with long execution times
	lock, err := clients.GetLock(pid, time.Minute)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	// TODO: Figure this out
	defer clients.ReleaseLock(lock)

	time.Sleep(time.Second * 5)

	offer := &models.Trade{}
	if err := clients.GetObject("offer", pid, offer); err == redis.Nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No offer found for this id"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// TODO: Connect to other services that do transactions
	now := time.Now()
	offer.Executed = &now

	models.DB.Create(offer)
	if offer.ID == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create trade"})
		return
	}

	if err := clients.DeleteObject("offer", pid); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"offer": offer})
}
