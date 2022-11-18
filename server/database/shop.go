package database

import (
	"context"
	"github.com/bberik/ecom-gin-react/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
)

func AddProductToShop(userID *string, productID *string, ctx context.Context) error {

	filter := bson.D{primitive.E{Key: "userId", Value: userID}}
	update := bson.D{{Key: "$push", Value: bson.D{primitive.E{Key: "products", Value: productID}}}}
	_, err := Shops.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func RemoveProductFromShop(userID *string, productID *string, ctx context.Context) error {
	filter := bson.D{primitive.E{Key: "userId", Value: userID}}
	update := bson.M{"$pull": bson.M{"products": productID}}
	_, err := Shops.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func GetAllShops(ctx context.Context) ([]models.Shop, error) {
	var shoplist []models.Shop
	cursor, err := Shops.Find(ctx, bson.D{{}})
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &shoplist)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return shoplist, nil
}
