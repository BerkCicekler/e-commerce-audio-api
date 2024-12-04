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

func (r *ProductRepo) FetchFeatured(searchText string, offset int64, requestData *model.ProductRequest) ([]model.Product, error) {
	var categoryFilter bson.M

	if requestData.CategoryIdHex != "" {
		categoryOBID, err := primitive.ObjectIDFromHex(requestData.CategoryIdHex)
		if err != nil {
			return nil, err
		}
		categoryFilter = bson.M{
			"category": categoryOBID,
		}

	}

	search := bson.M{
		"name": bson.M{
			"$regex":   searchText + ".*",
			"$options": "i", // Case-insensitive
		},
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

	var combinedFilter bson.M

	if categoryFilter == nil {
		combinedFilter = bson.M{
			"$and": []bson.M{
				search,
				priceFilter,
			},
		}
	} else {
		combinedFilter = bson.M{
			"$and": []bson.M{
				search,
				categoryFilter,
				priceFilter,
			},
		}
	}

	cur, err := r.MongoCollection.Find(context.TODO(), combinedFilter, requestData.RequestTOMongoDbOption(offset))
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
