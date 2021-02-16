package mongo

import (
	"context"
	"fmt"
	"github.com/birukbelay/item/entity"
	"github.com/birukbelay/item/packages/user"
	"github.com/birukbelay/item/utils/helpers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)
//var a =[]error{}
// UserMongoRepo implements the  interface
type UserMongoRepo struct {
	collection *mongo.Collection
}
//var ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
var ctx = context.Background()

// NewUserMongoRepo will create a new object of CategoriesGormRepo
func NewUserMongoRepo(C *mongo.Collection) user.UserRepository {
	return &UserMongoRepo{collection: C}
}
func (umr UserMongoRepo) GetUsers() ([]entity.User, []error) {
	var users []entity.User

	cursor, err := umr.collection.Find(ctx, bson.M{})

	if err != nil {
		var a []error
		a = append(a, err)
		return nil, a
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var usr entity.User
		cursor.Decode(&usr)
		users = append(users, usr)
	}
	if err := cursor.Err(); err != nil {
		var a []error
		a = append(a, err)
		return nil, a

	}
	return users, nil
}

func (umr UserMongoRepo) GetUser(id string) (*entity.User, []error) {
	usr := entity.User{}

	var a []error
	oid, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{{"_id", oid}}

	err := umr.collection.FindOne(ctx, filter).Decode(&usr)
	a = append(a, err)

	if err != nil {

		return nil, a
	}

	return &usr, nil
}

func (umr UserMongoRepo) UpdateUser(user *entity.User) (*entity.User, []error) {


	update := bson.M{"$set": user}
	filter := bson.D{{"_id", user.ID}}

	res, err := umr.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		var errs []error
		errs = append(errs, err)
		return nil, errs
	}

	fmt.Println(res.UpsertedCount)
	fmt.Println(res.UpsertedID)
	//fmt.Printf("%v, %T",res, res)
	return user, nil
}

func (umr UserMongoRepo) StoreUser(user *entity.User) (*entity.User, []error) {
	usr:=user
	//helpers.LogValue("reached store",user)
	id, err := umr.collection.InsertOne(ctx, usr)
	if err!=nil{
		fmt.Println("er-", err)
		helpers.LogValue("InsertOneError", err)

		var a []error
		a= append(a, err)
		return nil, a
	}

	usr.ID,_=id.InsertedID.(primitive.ObjectID)


	return usr, nil
}

//TODO change the return type
func (umr UserMongoRepo) DeleteUser(id string) (*entity.User, []error) {
	oid, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{{"_id", oid}}

	_, err := umr.collection.DeleteOne(ctx, filter)
	if err != nil {
		var errs []error
		errs = append(errs, err)
		return nil, errs
	}
	//fmt.Println(result.DeletedCount)
	return nil, nil
}



func (umr UserMongoRepo) UserByName(name string) (*entity.User, []error) {
	usr := entity.User{}

	var a []error

	filter := bson.D{{"username", name}}

	err := umr.collection.FindOne(ctx, filter).Decode(&usr)
	a = append(a, err)

	if err != nil {

		return nil, a
	}

	return &usr, nil
}

func (umr UserMongoRepo) UserByPhone(phone string) (*entity.User, []error) {
	usr := entity.User{}


	var a []error

	filter := bson.D{{"phone", phone}}

	err := umr.collection.FindOne(ctx, filter).Decode(&usr)
	a = append(a, err)

	if err != nil {

		return nil, a
	}

	return &usr, nil
}

func (umr UserMongoRepo) UserByEmail(email string) (*entity.User, []error) {
	usr := entity.User{}


	var a []error

	filter := bson.D{{"email", email}}

	err := umr.collection.FindOne(ctx, filter).Decode(&usr)
	a = append(a, err)

	if err != nil {
		return nil, a
	}

	return &usr, nil
}


func (umr UserMongoRepo) PhoneExists(phone string) bool {
	usr := entity.User{}
	filter := bson.D{{"phone", phone}}

	err := umr.collection.FindOne(ctx, filter).Decode(&usr)
	//fmt.Println(err)
	if err != nil {
		if err ==mongo.ErrNoDocuments{
			return false
		}
	}
	helpers.LogValue("PhoneExists FindOne Error", err)

	return true
}

func (umr UserMongoRepo) EmailExists(email string) bool {
	usr := entity.User{}

	filter := bson.D{{"email", email}}
	err := umr.collection.FindOne(ctx, filter).Decode(&usr)
	if err != nil {
		if err ==mongo.ErrNoDocuments{
			return false
		}
	}
	fmt.Println(usr)
	return true
}
















