package mongo

import (
	"context"
	//"context"
	"fmt"
	//"time"

	"go.mongodb.org/mongo-driver/bson"
	//"go.mongodb.org/mongo-driver/bson/primitive"
	//"go.mongodb.org/mongo-driver/mongo"

	"github.com/birukbelay/item/entity"

)
func (umr UserMongoRepo) AddSession(ctx context.Context, user *entity.User) (*entity.User, []error) {
	session := bson.D{
		{"$set", bson.D{
			{"session", user.Session},
		}},
	}
	//update := bson.M{"$set": user}

	filter := bson.D{{"_id", user.ID}}

	res, err := umr.collection.UpdateOne(ctx, filter, session)
	if err != nil {
		errs := []error{}
		errs = append(errs, err)
		return nil, errs
	}

	fmt.Println(res.UpsertedCount)
	fmt.Println(res.UpsertedID)
	//fmt.Printf("%v, %T",res, res)
	return user, nil
}
