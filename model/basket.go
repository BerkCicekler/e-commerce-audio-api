package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Basket struct {
	ID primitive.ObjectID			`bson:"_id"`
	UserID primitive.ObjectID		`bson:"userId"`
	ProductID primitive.ObjectID	`bson:"productId"`
	Count uint8						`bson:"count"`
}

type BasketResponseModel struct {
	ID primitive.ObjectID		`bson:"_id" json:"_id,omitempty"`
	Product Product				`bson:"product" json:"product"`
	Count uint8					`bson:"count" json:"count"`
}

type BasketUpdateRequestModel struct {
	ID primitive.ObjectID		`bson:"_id" json:"basketId"`
	Count int8					`bson:"count" json:"count"`
}

type BasketAddToBasketRequestModel struct {
	Product string				`json:"productId"`
}