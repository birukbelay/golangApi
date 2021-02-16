package mongo

import (
	"fmt"
	"github.com/birukbelay/item/entity"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strconv"
)

type CategoriesMongoRepo struct {
	collection *mongo.Collection
}


//TODO
func (cmr CategoriesMongoRepo) Categories(limit int , offset string) ([]entity.Categories, []error) {
	Categoriess :=[]entity.Categories{}

	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	ofsets, err := strconv.Atoi(offset)

	if err!=nil{
		var a []error
		a=append(a, err)
		return nil, a
	}


	page:= limit * (ofsets-1)
	fmt.Println("page", page)
	findOptions.SetSkip(int64(page))

	var errs []error
	cursor, err := cmr.collection.Find(ctx, bson.M{}, findOptions)

	if err != nil {
		errs = append(errs, err)
		return nil, errs
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var ctg entity.Categories
		cursor.Decode(&ctg)
		Categoriess = append(Categoriess, ctg)
	}
	if err := cursor.Err(); err != nil {
		errs = append(errs, err)
		return nil, errs
	}
	return Categoriess, nil
}

//TODO
func (cmr CategoriesMongoRepo) Category(id string) (*entity.Categories, []error) {
	categories := entity.Categories{}

	var a []error
	oid, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{{"_id", oid}}

	err := cmr.collection.FindOne(ctx, filter).Decode(&categories)
	a = append(a, err)

	if err != nil {

		return nil, a
	}

	return &categories, nil
}

//TODO
func (cmr CategoriesMongoRepo) UpdateCategories(categories *entity.Categories) (*entity.Categories, []error) {
	update := bson.M{"$set": categories}

	filter := bson.D{{"_id", categories.ID}}

	res, err := cmr.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		var errs []error
		errs = append(errs, err)
		return nil, errs
	}

	fmt.Println(res.UpsertedCount)
	fmt.Println(res.UpsertedID)
	//fmt.Printf("%v, %T",res, res)
	return categories, nil
}

//TODO
func (cmr CategoriesMongoRepo) DeleteCategories(id string) (*entity.Categories, []error) {
	oid, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{{"_id", oid}}

	_, err := cmr.collection.DeleteOne(ctx, filter)
	if err != nil {
		var errs []error
		errs = append(errs, err)
		return nil, errs
	}
	//fmt.Println(result.DeletedCount)
	return nil, nil
}
//TODO
func (cmr CategoriesMongoRepo) StoreCategories(categories *entity.Categories) (*entity.Categories, []error) {

	fmt.Println(categories)
	id, err := cmr.collection.InsertOne(ctx, categories)
	if err!=nil{
		var a []error
		a= append(a, err)
		return nil, a
	}

	categories.ID =id.InsertedID.(primitive.ObjectID)

	return categories, nil
}

func NewCategoriesMongoRepo(C *mongo.Collection) *CategoriesMongoRepo {
	return &CategoriesMongoRepo{collection: C}
}

