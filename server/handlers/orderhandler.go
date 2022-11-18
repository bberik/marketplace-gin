package handlers

import (
	"context"
	"errors"
	"github.com/bberik/ecom-gin-react/database"
	"github.com/bberik/ecom-gin-react/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"time"
)

type OrderHandler struct {
	orderCollection   *mongo.Collection
	productCollection *mongo.Collection
	userCollection    *mongo.Collection
	shopscollection   *mongo.Collection
}

func NewOrderHandler() *OrderHandler {
	return &OrderHandler{
		orderCollection:   database.Orders,
		productCollection: database.Products,
		userCollection:    database.Users,
		shopscollection:   database.Shops,
	}
}

func (h *OrderHandler) OrderDirect(c *gin.Context) {
	user_id := c.Query("userID")
	if user_id == "" {
		log.Println("user id is empty")
		_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
		return
	}

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var order models.Order

	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if validation_err := Validate.Struct(order); validation_err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validation_err.Error()})
		return
	}

	id := primitive.NewObjectID()
	order.OrderID = id.Hex()

	var orderstatus models.OrderUpdate
	orderstatus.Status = "Order Created"
	orderstatus.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.OrderStatus = make([]models.OrderUpdate, 0)
	order.OrderStatus = append(order.OrderStatus, orderstatus)

	filter := bson.D{primitive.E{Key: "userID", Value: user_id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "orders", Value: order}}}}
	result, err := h.userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	if result.MatchedCount == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "UserID not found"})
		return
	}

	_, err = h.orderCollection.InsertOne(ctx, order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	filter = bson.D{primitive.E{Key: "products", Value: order.OrderContent.ProductID}}
	update = bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "orders", Value: order.OrderID}}}}
	_, err = h.shopscollection.UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	defer cancel()
	c.IndentedJSON(http.StatusOK, "Order successfully processed.")
}

func (h *OrderHandler) OrderFromCart(c *gin.Context) {
	user_id := c.Query("userID")
	if user_id == "" {
		log.Println("user id is empty")
		_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
		return
	}

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var user models.User
	var address models.Address

	if err := c.ShouldBindJSON(&address); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.userCollection.FindOne(ctx, bson.M{"userID": user_id}).Decode(&user)
	defer cancel()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "UserID is not a valid"})
		return
	}

	for _, content := range user.Cart {
		// orders = append(orders, models.Order{})
		var order models.Order
		id := primitive.NewObjectID()
		order.OrderID = id.Hex()
		order.OrderContent.ProductID = content.ProductID
		order.OrderContent.Quantity = content.Quantity
		order.OrderContent.ItemID = content.ItemID
		order.ShippingAddress = address
		var orderstatus models.OrderUpdate
		orderstatus.Status = "Order Created"
		orderstatus.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		order.OrderStatus = make([]models.OrderUpdate, 0)
		order.OrderStatus = append(order.OrderStatus, orderstatus)

		_, err = h.orderCollection.InsertOne(ctx, order)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		filter := bson.D{primitive.E{Key: "userID", Value: user.UserID}}
		update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "orders", Value: order}}}}
		_, err = h.userCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		filter = bson.D{primitive.E{Key: "products", Value: order.OrderContent.ProductID}}
		update = bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "orders", Value: order.OrderID}}}}
		_, err = h.shopscollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
	}

	emptycart := make([]models.Content, 0)
	filter := bson.D{primitive.E{Key: "userID", Value: user.UserID}}
	update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "cart", Value: emptycart}}}}
	_, err = h.userCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.IndentedJSON(http.StatusOK, "Orders from cart successfully processed!")
}
