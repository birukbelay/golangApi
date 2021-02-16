package mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"

	"github.com/birukbelay/item/entity"
	"github.com/birukbelay/item/models/items"
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

var ctx, _ = context.WithTimeout(context.Background(), 20*time.Second)
//var ctx = context.Background()

func (pmr ProductMongoRepo) Items(limit, offset int) ([]entity.Item, []error) {

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
func (pmr ProductMongoRepo) Item(id string) (*entity.Item, []error) {
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
func (pmr ProductMongoRepo) StoreItem(item *entity.Item) (*entity.Item, []error) {

	itm := item

	_, err := pmr.collection.InsertOne(ctx, itm)
	fmt.Println(",,,,,,,,.")
	if err!=nil{
		fmt.Println("err..",err, "ctX",ctx)

		a:=[]error{}
		a= append(a, err)
		return nil, a
	}

	return itm, nil
}

// UpdateItem ...
func (pmr ProductMongoRepo) UpdateItem(item *entity.Item) (*entity.Item, []error) {
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
	//fmt.Printf("%v, %T",res, res)
	return item, nil
}

// DeleteItem ...
func (pmr ProductMongoRepo) DeleteItem(id string) (int64, []error) {
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

func (pmr ProductMongoRepo) ItemsByCategories(limit, offset int, categories string) ([]entity.Item, []error) {

	itm := []entity.Item{}

	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))

	var errs []error
	cursor, err := pmr.collection.Find(ctx, bson.M{"categories": categories}, findOptions)

	if err != nil {

		errs = append(errs, err)
		return nil, errs
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var item entity.Item
		err := cursor.Decode(&item)
		if err != nil {
			errs = append(errs, err)
			return nil, errs
		}
		itm = append(itm, item)
	}

	if err := cursor.Err(); err != nil {
		errs = append(errs, err)
		return nil, errs
	}

	return itm, nil

}
//TODO
func (pmr ProductMongoRepo) StoreManyItems(items []interface{}) ([]interface{}, []error) {

	//items := []interface{}{item1, item2}

	insertManyResult, err := pmr.collection.InsertMany(ctx, items)
	if err != nil {
		return nil , nil
	}


	return insertManyResult.InsertedIDs, nil
}
