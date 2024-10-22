package repository

import (
	"context"
	"log"

	"github.com/BerkCicekler/e-commerce-audio-api/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CategoriesRepo struct {
	MongoCollection *mongo.Collection
}

func (r *CategoriesRepo) GetCategories() (*[]model.Category, error) {
	cur, err := r.MongoCollection.Find(context.TODO(), bson.M{})
	defer cur.Close(context.TODO())
	if err != nil {
		return nil, err
	}
	
	var results []model.Category
	for (cur.Next(context.TODO())) {
				var elem model.Category
				err := cur.Decode(&elem)
				if err != nil {
					log.Fatal(err)
				}
		
				results =append(results, elem)
	}
	if err := cur.Err(); err != nil {
        log.Fatal(err)
    }
	return &results, nil

}