package repository

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/BerkCicekler/e-commerce-audio-api/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BasketRepo struct {
	MongoCollection *mongo.Collection
}

func (r *BasketRepo) FetchUserBasket(uId *primitive.ObjectID) (*[]model.BasketResponseModel, error) {
	filter := bson.M{
        "userId": uId,
    }

	pipeline := mongo.Pipeline{
        {{"$match", filter}},
        {{"$lookup", bson.M{
            "from":         "products",
            "localField":   "productId",
            "foreignField": "_id",
            "as":           "product",
        }}},
		{{"$unwind", "$product"}},
    }

	cur, err := r.MongoCollection.Aggregate(context.TODO(), pipeline)
	if err != nil {
        log.Fatal(err)
    }
    defer cur.Close(context.Background())
    results := []model.BasketResponseModel{}


    if err = cur.All(context.TODO(), &results); err != nil {
        log.Fatal(err)
		return nil, nil
    }
	return &results, nil
	
}

func (r *BasketRepo) AddToBasket(productId, uId *primitive.ObjectID) (primitive.ObjectID, error) {
	isExist, err := r.isProductAlreadyInBasket(productId, uId)
	if err != nil  {
		return primitive.NilObjectID, err
	}
	if isExist {
		return primitive.NilObjectID, errors.New("Item already in basket")
	}
	obId := primitive.NewObjectID()
	
	fmt.Println(obId)
	_, err = r.MongoCollection.InsertOne(context.TODO(), model.Basket{
		ID: obId,
		UserID: *uId,
		ProductID: *productId,
		Count: 1,
	})
	if err != nil {
		return primitive.NilObjectID, err
	}
	
	return primitive.NilObjectID, nil
}

func (r *BasketRepo) AddToItemCount(value int8, BasketID *primitive.ObjectID) error {
	filter := bson.M{
		"_id": BasketID,
	}

	update := bson.M{
		"$inc": bson.M{"count": value},
	}

	result, err := r.MongoCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("The product couldn't found in basket")
	}

	return nil
}

func (r *BasketRepo) DeleteItem(basketId *primitive.ObjectID) error {
	_, err := r.MongoCollection.DeleteOne(context.TODO(), bson.M{
		"_id": basketId,
	})
	return err
}

func (r *BasketRepo) ClearUserBasket(userId *primitive.ObjectID) error {
	_, err := r.MongoCollection.DeleteMany(context.TODO(), bson.M{
		"userId": userId,
	})
	return err
}

func (r *BasketRepo) isProductAlreadyInBasket(productId,userId *primitive.ObjectID) (bool, error) {
	result, err := r.MongoCollection.CountDocuments(context.TODO(), bson.M{
		"userId": userId,
		"productId": productId,
	})
	return result > 0, err
}