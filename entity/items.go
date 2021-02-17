package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Item struct{
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string `bson:"name" json:"name"`

	Description string   `bson:"description,omitempty" json:"description"`
	Image       string   `bson:"image,omitempty" json:"image"`
	Categories  []string `bson:"categories,omitempty" json:"categories"`
	Price       int      `bson:"year,omitempty" json:"price"`

	Type        string   `bson:"type,omitempty" json:"type"`


}


type Categories struct{
	ID primitive.ObjectID `bson:"_id" json:"id`
	Name string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
	Image string `json:"image" bson:"image`
}

