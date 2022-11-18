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

type ShopHandler struct {
	orderCollection   *mongo.Collection
	productCollection *mongo.Collection
	userCollection    *mongo.Collection
}

func NewShopHandler() *ShopHandler {
	return &ShopHandler{
		orderCollection:   database.Orders,
		productCollection: database.Products,
		userCollection:    database.Users,
	}
}

func (h *ShopHandler) CreateProduct(c *gin.Context) {

	userQueryID := c.Query("userID")
	if userQueryID == "" {
		log.Println("user id is empty")
		_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
		return
	}

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var product models.Product

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if validation_err := Validate.Struct(product); validation_err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validation_err.Error()})
		return
	}

	id := primitive.NewObjectID()
	product.ProductID = id.Hex()
	for _, v := range product.Items {
		iid := primitive.NewObjectID()
		v.ItemID = iid.Hex()
	}

	err := database.CreateProduct(&product, ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	err = database.AddProductToShop(&userQueryID, &product.ProductID, ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	defer cancel()

	c.IndentedJSON(http.StatusOK, "Product added successfully")
}

func (h *ShopHandler) UpdateProduct(c *gin.Context) {
	userQueryID := c.Query("userID")
	if userQueryID == "" {
		log.Println("user id is empty")
		_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
		return
	}

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var product models.Product

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if validation_err := Validate.Struct(product); validation_err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validation_err.Error()})
		return
	}

	for _, v := range product.Items {
		if v.ItemID == "" {
			iid := primitive.NewObjectID()
			v.ItemID = iid.Hex()
		}
	}

	err := database.UpdateProduct(&product, ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	defer cancel()

	c.IndentedJSON(http.StatusOK, "Product updated successfully")
}

func (h *ShopHandler) DeleteProduct(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	productQueryID := c.Query("pid")
	if productQueryID == "" {
		log.Println("product id is empty")
		_ = c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
		return
	}
	userQueryID := c.Query("userID")
	if userQueryID == "" {
		log.Println("user id is empty")
		_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
		return
	}

	err := database.DeleteProduct(&productQueryID, ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	err = database.RemoveProductFromShop(&userQueryID, &productQueryID, ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
	}

	defer cancel()

	c.IndentedJSON(http.StatusOK, "Product deleted successfully")

}

func (h *ShopHandler) UpdateOrderStatus(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	orderQueryID := c.Query("oid")
	if orderQueryID == "" {
		log.Println("order id is empty")
		_ = c.AbortWithError(http.StatusBadRequest, errors.New("order id is empty"))
		return
	}
	userQueryID := c.Query("userID")
	if userQueryID == "" {
		log.Println("user id is empty")
		_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
		return
	}

	count, err := database.Shops.CountDocuments(ctx, bson.M{"userId": userQueryID, "orders": orderQueryID})
	if err != nil {
		log.Panic(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	if count != 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order is not for this shop."})
		return
	}

	var orderupdate *models.OrderUpdate
	if err := c.ShouldBindJSON(&orderupdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	orderupdate.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	filter := bson.D{primitive.E{Key: "orderID", Value: orderQueryID}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "orderStatus", Value: orderupdate}}}}
	_, err = h.orderCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	defer cancel()
	c.IndentedJSON(200, "Order Status Updated Successfully")
}

func (h *ShopHandler) GetAllShops(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	shoplist, err := database.GetAllShops(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	defer cancel()
	c.IndentedJSON(200, shoplist)
}
