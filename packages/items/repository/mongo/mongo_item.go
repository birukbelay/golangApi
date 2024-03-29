package mongo

import (
	"context"
	"errors"
	"fmt"
	"github.com/birukbelay/item/entity"
	"github.com/birukbelay/item/packages/items"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)
//var a =[]error{}
// UserMongoRepo implements the items.CategoriesRepository interface
type ProductMongoRepo struct {
	collection *mongo.Collection

}

// NewUserMongoRepo will create a new object of CategoriesGormRepo
func NewProductMongoRepo(C *mongo.Collection) items.ItemRepository {
	return &ProductMongoRepo{collection: C}
}


func (pmr ProductMongoRepo) Items(ctx context.Context, limit, offset int) ([]entity.Item, []error) {

	itm :=[]entity.Item{}

	cursor, err := pmr.collection.Find(ctx, bson.M{})

	if err != nil {
		a := []error{}
		a = append(a, err)
		return nil, a
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var item entity.Item
		cursor.Decode(&item)
		itm = append(itm, item)
	}
	if err := cursor.Err(); err != nil {
		a := []error{}
		a = append(a, err)
		return nil, a

	}
	return itm, nil


}

// Item gets a single items from database
func (pmr ProductMongoRepo) Item(ctx context.Context, id string) (*entity.Item, []error) {
	item := entity.Item{}


	a := []error{}
	oid, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{{"_id", oid}}

	err := pmr.collection.FindOne(ctx, filter).Decode(&item)
	a = append(a, err)

	if err != nil {

		return nil, a
	}

	return &item, nil
}

// StoreItem stores an items into database
func (pmr ProductMongoRepo) StoreItem(ctx context.Context, item *entity.Item) (*entity.Item, []error) {


	ids:= primitive.NewObjectID()
	item.ID=ids

	id, err := pmr.collection.InsertOne(ctx, item)
	fmt.Println(",,,,,,,,.")
	if err!=nil{
		fmt.Println("err..",err, "ctX",ctx)

		a:=[]error{}
		a= append(a, err)
		return nil, a
	}
	item.ID=id.InsertedID.(primitive.ObjectID)

	return item, nil
}

// UpdateItem ...
func (pmr ProductMongoRepo) UpdateItem(ctx context.Context, item *entity.Item) (*entity.Item, []error) {
	update := bson.M{"$set": item}

	filter := bson.D{{"_id", item.ID}}

	res, err := pmr.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		errs := []error{}
		errs = append(errs, err)
		return nil, errs
	}

	fmt.Println(res.UpsertedCount)
	fmt.Println(res.UpsertedID)
	if res.MatchedCount>0{
		//fmt.Printf("%v, %T",res, res)
		return item, nil

	}
	var errs []error
	errs=append(errs, errors.New("not updated"))
	return nil, errs

}

// DeleteItem ...
func (pmr ProductMongoRepo) DeleteItem(ctx context.Context, id string) (int64, []error) {
	oid, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{{"_id", oid}}

	result, err := pmr.collection.DeleteOne(ctx, filter)
	if err != nil {
		errs := []error{}
		errs = append(errs, err)
		return 0, errs
	}
	//fmt.Println(result.DeletedCount)
	return result.DeletedCount, nil
}
