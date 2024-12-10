package repository

import (
	"context"

	"github.com/BerkCicekler/e-commerce-audio-api/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UsersRepo struct {
	MongoCollection *mongo.Collection
}

func (r *UsersRepo) InsertUser(user *model.User) (*mongo.InsertOneResult, error) {
	result, err := r.MongoCollection.InsertOne(context.Background(), user)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *UsersRepo) InsertOAuthUser(user *model.OAuthUser) (*mongo.InsertOneResult, error) {
	result, err := r.MongoCollection.InsertOne(context.Background(), user)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *UsersRepo) FindUserByEmail(email string) (*model.User, error) {
	var user model.User
	filter := bson.M{"email": email}
	err := r.MongoCollection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UsersRepo) FindUserById(Id primitive.ObjectID) (*model.User, error) {
	var user model.User
	filter := bson.M{"_id": Id}
	err := r.MongoCollection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UsersRepo) UpdateUserById(user *model.User) (int64, error) {
	result, err := r.MongoCollection.UpdateOne(context.Background(),
		bson.D{{Key: "_id", Value: user.ID}},
		bson.D{{Key: "$set", Value: user}})

	if err != nil {
		return 0, err
	}
	return result.ModifiedCount, nil
}
