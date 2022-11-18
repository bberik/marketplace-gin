package database

import (
	"context"
	"errors"
	"github.com/bberik/ecom-gin-react/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateUser(user *models.User, ctx context.Context) error {
	_, err := Users.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func GetUser(userID *string, ctx context.Context) (*models.User, error) {
	var user *models.User
	query := bson.D{bson.E{Key: "userID", Value: userID}}
	err := Users.FindOne(ctx, query).Decode(&user)
	return user, err
}

func UpdateUser(user *models.User, ctx context.Context) error {

	filter := bson.D{primitive.E{Key: "userID", Value: user.UserID}}
	update := bson.D{primitive.E{Key: "$set", Value: bson.D{
		primitive.E{Key: "username", Value: user.Username},
		primitive.E{Key: "fullName", Value: user.FullName}}}}
	result, _ := Users.UpdateOne(ctx, filter, update)
	if result.MatchedCount != 1 {
		return errors.New("no matched user found for update")
	}
	return nil
}

func DeleteUser(userID *string, ctx context.Context) error {
	filter := bson.D{primitive.E{Key: "userID", Value: userID}}
	result, _ := Users.DeleteOne(ctx, filter)
	if result.DeletedCount != 1 {
		return errors.New("no matched user found for delete")
	}
	return nil
}
