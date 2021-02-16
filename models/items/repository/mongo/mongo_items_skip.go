package mongo

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"reflect"
	"strconv"
	//"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/birukbelay/item/entity"
)



//ItemsBySkipFilter ...
func (pmr ProductMongoRepo) ItemsByFilter(limit int, offsetCursor, searchField,
	categories, brand, types string,
	sort string, sortWay int,
	minYear, maxYear int) ([]entity.Item,string, string, []error) {

	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))

	//sort_way is 1 or -1
	if sort ==""{
		sort ="Featured"
	}
	findOptions.SetSort(bson.D{{sort, sortWay},{"_id",1}})
	offsetsVal, err := strconv.Atoi(offsetCursor)

	page:= limit * (offsetsVal-1)
	fmt.Println("page", page)
	findOptions.SetSkip(int64(page))
	query := bson.D{}
	//a:=bson.D{{"brand", bson.D{{ "$in", bson.A{"samsung", "Bob"} }}  }}
	
	// bson.M{"categoriesId": categoriesId, "brandId": brandId, "types": types, "$and": AndQuery}
	// bson.D{{"categoriesId": categoriesId}, "brandId": brandId, "types": types, "$and": AndQuery}

	//------------------------------------- 4 query Filters -----------------------------------

	if categories !="" && categories !="undefined"{
		fmt.Println(reflect.TypeOf(categories))
		query = append(query, bson.E{"categories", categories})
	}

	if brand !=""{
		query = append(query, bson.E{"brand", brand})
	}

	if types !=""{
		query = append(query, bson.E{"type", types})
	}


	if minYear !=0 || maxYear !=0{
		var yearQuery = bson.D{}

		if maxYear != 0{
			yearQuery =append(yearQuery, bson.E{Key: "$lte", Value: maxYear})
		}
		if minYear !=0{
			yearQuery =append(yearQuery, bson.E{"$gte", minYear})
		}
		query = append(query, bson.E{"year", yearQuery})

	}

	//--------------  -------------------------------------------////-----------------------

	if searchField !=""{
		var searchQuery= bson.M{"$search":searchField}
		// query["$text"]= searchQuery
		query = append(query, bson.E{"$text",searchQuery})
	}


	var errs []error
	cursor, err := pmr.collection.Find(ctx, query, findOptions)

	fmt.Println("query-", query)

	if err != nil {

		errs = append(errs, err)

		return nil,"","", errs
	}
	defer cursor.Close(ctx)

	//var ofsetValue string
	i:=0
	var items []entity.Item
	for cursor.Next(ctx) {
		var item entity.Item
		err := cursor.Decode(&item)
		i++

		if err != nil {
			errs = append(errs, err)
			return nil, "","",errs
		}
		items = append(items, item)



	}

	if err := cursor.Err(); err != nil {
		errs = append(errs, err)
		return nil,"","", errs
	}
	fmt.Println("items--", items)
	goPrev := fmt.Sprintf("%v_%v", "firstValue", "stringFirstValue")
	goNext := fmt.Sprintf("%v_%v", "lastValue", "stringLastValue")
	return items, goPrev, goNext, nil


}

