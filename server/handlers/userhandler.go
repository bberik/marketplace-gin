package handlers

import (
	"context"
	"errors"
	"github.com/bberik/ecom-gin-react/database"
	"github.com/bberik/ecom-gin-react/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

var Validate = validator.New()

type UserHandler struct {
	userCollection *mongo.Collection
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		userCollection: database.Users,
	}
}

func hashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func verifyPassword(userPassword string, inputPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(inputPassword), []byte(userPassword))
	valid := true
	msg := ""
	if err != nil {
		msg = "Login Or Password is Incorrect"
		valid = false
	}
	return valid, msg
}

func (h *UserHandler) SignUp(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if validation_err := Validate.Struct(user); validation_err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validation_err.Error()})
		return
	}

	count, err := h.userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
	if err != nil {
		log.Panic(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
		return
	}

	password := hashPassword(*user.Password)
	user.Password = &password

	id := primitive.NewObjectID()
	user.UserID = id.Hex()
	user.Cart = make([]models.Content, 0)
	user.UserAddress = make([]models.Address, 0)
	user.Orders = make([]models.Order, 0)
	dberr := database.CreateUser(&user, ctx)
	if dberr != nil {
		c.JSON(http.StatusBadGateway, gin.H{"message": dberr.Error()})
		return
	}
	defer cancel()
	c.JSON(http.StatusOK, gin.H{"message": "Successfully SignedUp!"})
}

func (h *UserHandler) SignInUser(c *gin.Context) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var user models.User
	var foundUser models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
	defer cancel()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	passwordIsCorrect, msg := verifyPassword(*user.Password, *foundUser.Password)

	if !passwordIsCorrect {
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		return
	}
	token := GenerateToken(foundUser.UserID)
	c.JSON(http.StatusOK, gin.H{"msg": "Successfully SignIN", "token": token})
}

func (h *UserHandler) GetUser(c *gin.Context) {
	user_id := c.Query("userID")
	if user_id == "" {
		log.Println("user id is empty")
		_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	user, err := database.GetUser(&user_id, ctx)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	defer cancel()
	c.IndentedJSON(http.StatusOK, user)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if validation_err := Validate.Struct(user); validation_err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validation_err.Error()})
		return
	}

	err := database.UpdateUser(&user, ctx)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
	}
	defer cancel()
	c.IndentedJSON(http.StatusOK, "User Information successfully updated.")
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	user_id := c.Query("userID")
	if user_id == "" {
		log.Println("user id is empty")
		_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := database.DeleteUser(&user_id, ctx)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	defer cancel()
	c.IndentedJSON(http.StatusOK, "User successfully deleted")
}

func (h *UserHandler) OpenShop(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var newshop models.Shop
	if err := c.ShouldBindJSON(&newshop); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if validation_err := Validate.Struct(newshop); validation_err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validation_err.Error()})
		return
	}
	var user models.User
	err := h.userCollection.FindOne(ctx, bson.M{"userID": newshop.UserID}).Decode(&user)
	defer cancel()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	count, err := database.Shops.CountDocuments(ctx, bson.M{"userId": newshop.UserID})
	if err != nil {
		log.Panic(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User Already has a shop"})
		return
	}

	newshop.Products = make([]*string, 0)
	newshop.Orders = make([]*string, 0)

	_, err = database.Shops.InsertOne(ctx, newshop)
	if err != nil {
		c.IndentedJSON(500, "something Went wrong")
		return
	}
	defer cancel()
	ctx.Done()
	c.IndentedJSON(200, "Congratulations! You opened your own shop!")
}

func (h *UserHandler) AddAddress(c *gin.Context) {
	user_id := c.Query("userID")
	if user_id == "" {
		log.Println("user id is empty")
		_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
		return
	}

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var address models.Address

	if err := c.ShouldBindJSON(&address); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if validation_err := Validate.Struct(address); validation_err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validation_err.Error()})
		return
	}

	filter := bson.D{primitive.E{Key: "userID", Value: user_id}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "userAddress", Value: address}}}}
	result, err := h.userCollection.UpdateOne(ctx, filter, update)
	if result.MatchedCount == 0 {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "There is no user with given ID"})
		return
	}
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Cannot add address to user"})
		return
	}

	defer cancel()

	c.IndentedJSON(http.StatusOK, "Address added to the user!")
}

func (h *UserHandler) UpdateAddress(c *gin.Context) {
	user_id := c.Query("userID")
	if user_id == "" {
		log.Println("user id is empty")
		_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
		return
	}

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var addresses []models.Address

	if err := c.ShouldBindJSON(&addresses); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for _, address := range addresses {
		if validation_err := Validate.Struct(address); validation_err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validation_err.Error()})
			return
		}
	}

	filter := bson.D{primitive.E{Key: "userID", Value: user_id}}
	update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "userAddress", Value: addresses}}}}
	result, err := h.userCollection.UpdateOne(ctx, filter, update)
	if result.MatchedCount == 0 {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "There is no user with given ID"})
		return
	}
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Cannot update user address"})
		return
	}
	defer cancel()

	c.IndentedJSON(http.StatusOK, "Updated user's addresses!")
}
