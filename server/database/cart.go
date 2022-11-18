package database

import (
	"context"
	"github.com/bberik/ecom-gin-react/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func AddToCart(content *models.Content, userID *string, ctx context.Context) error {
	filter := bson.D{primitive.E{Key: "userID", Value: userID}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "cart", Value: content}}}}
	_, err := Users.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func UpdateCart(cart []models.Content, userID *string, ctx context.Context) error {
	filter := bson.D{primitive.E{Key: "userID", Value: userID}}
	update := bson.D{{Key: "$set", Value: bson.D{primitive.E{Key: "cart", Value: cart}}}}
	_, err := Users.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}
