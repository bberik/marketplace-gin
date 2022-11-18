package database

import (
	"context"
	"errors"
	"github.com/bberik/ecom-gin-react/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
)

func CreateProduct(product *models.Product, ctx context.Context) error {
	_, err := Products.InsertOne(ctx, product)
	if err != nil {
		return err
	}
	return nil
}

func GetAllProducts(ctx context.Context) ([]models.Product, error) {
	var productlist []models.Product
	cursor, err := Products.Find(ctx, bson.D{{}})
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &productlist)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return productlist, nil
}

func GetProduct(productID *string, ctx context.Context) (models.Product, error) {
	var product models.Product
	query := bson.D{bson.E{Key: "pID", Value: productID}}
	err := Users.FindOne(ctx, query).Decode(&product)
	return product, err
}

func GetProductByCategory(category *string, ctx context.Context) ([]models.Product, error) {
	var productlist []models.Product
	filter := bson.D{primitive.E{Key: "category", Value: category}}
	cursor, err := Products.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	err = cursor.All(ctx, &productlist)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return productlist, nil
}

func SearchProduct(regex *string, ctx context.Context) ([]models.Product, error) {
	var searchproducts []models.Product
	searchquerydb, err := Products.Find(ctx, bson.M{"product_name": bson.M{"$regex": regex}})
	if err != nil {
		return nil, err
	}
	err = searchquerydb.All(ctx, &searchproducts)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer searchquerydb.Close(ctx)

	if err := searchquerydb.Err(); err != nil {
		log.Println(err)
		return nil, err
	}

	return searchproducts, nil
}

func UpdateProduct(product *models.Product, ctx context.Context) error {

	filter := bson.D{primitive.E{Key: "pID", Value: product.ProductID}}
	update := bson.D{primitive.E{Key: "$set", Value: bson.D{
		primitive.E{Key: "productName", Value: product.ProductName},
		primitive.E{Key: "description", Value: product.Description},
		primitive.E{Key: "image", Value: product.Images},
		primitive.E{Key: "items", Value: product.Items},
		primitive.E{Key: "category", Value: product.Category}}}}

	result, _ := Products.UpdateOne(ctx, filter, update)
	if result.MatchedCount != 1 {
		return errors.New("no matched user found for update")
	}
	return nil
}

func DeleteProduct(productID *string, ctx context.Context) error {
	filter := bson.D{primitive.E{Key: "pID", Value: productID}}
	result, _ := Products.DeleteOne(ctx, filter)
	if result.DeletedCount != 1 {
		return errors.New("no matched product found for delete")
	}
	return nil
}
