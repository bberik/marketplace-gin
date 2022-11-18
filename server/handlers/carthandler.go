package handlers

import (
	"context"
	"errors"
	"github.com/bberik/ecom-gin-react/database"
	"github.com/bberik/ecom-gin-react/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"time"
)

type CartHandler struct {
	productCollection *mongo.Collection
	userCollection    *mongo.Collection
}

func NewCartHandler() *CartHandler {
	return &CartHandler{
		productCollection: database.Products,
		userCollection:    database.Users,
	}
}

func (h *CartHandler) AddToCart(c *gin.Context) {
	user_id := c.Query("userID")
	if user_id == "" {
		log.Println("user id is empty")
		_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
		return
	}

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var user models.User
	var content models.Content

	if err := c.ShouldBindJSON(&content); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.userCollection.FindOne(ctx, bson.M{"userID": user_id}).Decode(&user)
	defer cancel()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "UserID is not a valid"})
		return
	}

	err = database.AddToCart(&content, &user.UserID, ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not add product to cart of user"})
		return
	}

	defer cancel()

	c.IndentedJSON(http.StatusOK, "Product added to the cart of user!")
}

func (h *CartHandler) UpdateCart(c *gin.Context) {
	user_id := c.Query("userID")
	if user_id == "" {
		log.Println("user id is empty")
		_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
		return
	}

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var user models.User
	var cart []models.Content

	if err := c.ShouldBindJSON(&cart); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.userCollection.FindOne(ctx, bson.M{"userID": user_id}).Decode(&user)
	defer cancel()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "UserID is not a valid"})
		return
	}

	err = database.UpdateCart(cart, &user.UserID, ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update cart of user"})
		return
	}

	defer cancel()

	c.IndentedJSON(http.StatusOK, "Updated user's cart!")
}
