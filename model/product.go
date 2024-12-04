package model

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Product struct {
	ID          primitive.ObjectID `bson:"_id" json:"_id,omitempty"`
	PictureName string             `json:"pictureName" bson:"pictureName"`
	Name        string             `json:"name" bson:"name"`
	Price       float32            `json:"price" bson:"price"`
}

type ProductRequest struct {
	CategoryIdHex string `json:"category"`
	StartIndex    uint32 `json:"startIndex"`
	Search        string `json:"search"`
	SortBy        string `json:"sortBy"`
	MinPrice      uint32 `json:"minPrice"`
	MaxPrice      uint32 `json:"maxPrice"`
}

func (m *ProductRequest) RequestTOMongoDbOption(offset int64) *options.FindOptions {
	sortOrder := 0
	switch m.SortBy {
	case "desc":
		sortOrder = -1
	case "asc":
		sortOrder = 1
	}
	findOptions := options.Find()
	if sortOrder != 0 {
		findOptions.SetSort(bson.D{{Key: "price", Value: sortOrder}})
	}
	findOptions.SetSkip(offset)
	findOptions.SetLimit(10)
	return findOptions
}
