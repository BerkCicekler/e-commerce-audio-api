package repository

import (
	"context"
	"log"

	"github.com/BerkCicekler/e-commerce-audio-api/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductRepo struct {
	MongoCollection *mongo.Collection
}

func (r *ProductRepo) FetchFeatured(requestData *model.ProductRequest) ([]model.Product, error){
	categoryOBID, err:=  primitive.ObjectIDFromHex(requestData.CategoryIdHex)
	if err != nil{
		return nil, err
	}
    
    search := bson.M{
        "name": bson.M{
            "$regex":   requestData.Search + ".*", 
            "$options": "i",                       // Case-insensitive
        },
    }

    filter := bson.M{
        "category": categoryOBID,
    }

	priceFilter := bson.M{}
    if requestData.MinPrice > 0 {
        priceFilter["price"] = bson.M{"$gte": requestData.MinPrice}
    }
    if requestData.MaxPrice > 0 {
        if _, exists := priceFilter["price"]; exists {
            priceFilter["price"].(bson.M)["$lte"] = requestData.MaxPrice
        } else {
            priceFilter["price"] = bson.M{"$lte": requestData.MaxPrice}
        }
    }

    combinedFilter := bson.M{
        "$and": []bson.M{
            search,
            filter,
			priceFilter,
        },
    }
	cur, err := r.MongoCollection.Find(context.TODO(), combinedFilter, requestData.RequestTOMongoDbOption())
	if err != nil {
        log.Fatal(err)
    }
    defer cur.Close(context.Background())
    results := []model.Product{}


    if err = cur.All(context.TODO(), &results); err != nil {
        log.Fatal(err)
		return nil, nil
    }
	return results, nil
}