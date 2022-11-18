package handlers

import (
	"context"
	"errors"
	"github.com/bberik/ecom-gin-react/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"time"
)

type ProductHandler struct {
	productCollection *mongo.Collection
	userCollection    *mongo.Collection
	shopCollection    *mongo.Collection
}

func NewProductHandler() *ProductHandler {
	return &ProductHandler{
		productCollection: database.Products,
		userCollection:    database.Users,
		shopCollection:    database.Shops,
	}
}

func (h *ProductHandler) GetAllProducts(c *gin.Context) {

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	productlist, err := database.GetAllProducts(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	defer cancel()
	c.IndentedJSON(200, productlist)
}

func (h *ProductHandler) GetProduct(c *gin.Context) {
	product_id := c.Query("pID")
	if product_id == "" {
		log.Println("product id is empty")
		_ = c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	product, err := database.GetProduct(&product_id, ctx)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	defer cancel()
	c.IndentedJSON(http.StatusOK, product)
}

func (h *ProductHandler) GetProductsByCategory(c *gin.Context) {
	category := c.Query("category")
	if category == "" {
		log.Println("Category is empty")
		_ = c.AbortWithError(http.StatusBadRequest, errors.New("category is empty"))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	products, err := database.GetProductByCategory(&category, ctx)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	defer cancel()
	c.IndentedJSON(http.StatusOK, products)
}

func (h *ProductHandler) SearchProduct(c *gin.Context) {
	queryParam := c.Query("name")
	if queryParam == "" {
		log.Println("query is empty")
		c.Header("Content-Type", "application/json")
		c.JSON(http.StatusNotFound, gin.H{"Error": "Invalid Search Index"})
		c.Abort()
		return
	}
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	products, err := database.SearchProduct(&queryParam, ctx)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	defer cancel()
	c.IndentedJSON(http.StatusOK, products)
}
