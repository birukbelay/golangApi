package entity

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Item struct{
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string `bson:"name" json:"name"`

	Description string   `bson:"description,omitempty" json:"description"`
	Image       string   `bson:"image,omitempty" json:"image"`
	Categories  []string `bson:"categories,omitempty" json:"categories"`
	Price       int      `bson:"year,omitempty"`
	Language    string   `bson:"language,omitempty"`
	Type        string   `bson:"type,omitempty"`


}


type Categories struct{
	ID primitive.ObjectID
	Name string `json:"name"`
	Description string `json:"description"`
	Image string `json:"image"`
}

