package productHandler

import (
	"encoding/json"
	"fmt"
	"github.com/birukbelay/item/utils/global"
	"github.com/birukbelay/item/utils/helpers"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"strconv"
)

// GetItems handles GET /v1/admin/Items request
func (aih *AdminItemHandler) GetFilteredItems(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	query := r.URL.Query()
	fmt.Println("URLQuery", query)

	categories :=query.Get("categories")
	brand :=query.Get("brand")
	types :=query.Get("types")
	sort :=query.Get("sort")


	offsetValue :=query.Get("offsetValue")
	searchField :=query.Get("searchField")



	limit, err := strconv.Atoi(query.Get("limit"))
	if err != nil {
		helpers.HandleErr(w, err, global.StatusBadRequest, 400)
		limit = global.MaxItemLimit
	}
	if limit > global.MaxItemLimit {
		limit = global.MaxItemLimit
	}


	sortWay, err := strconv.Atoi(query.Get("sort_way"))
	if err != nil {
		//helpers.HandleErr(w, err, global.StatusBadRequest, 400)
		sortWay = -1

	}

	minPrice, err := strconv.Atoi(query.Get("min_price"))
	if err != nil {
		//helpers.HandleErr(w, err, global.StatusBadRequest, 400)
		minPrice = 0

	}

	maxPrice, err := strconv.Atoi(query.Get("max_price"))
	if err != nil {
		//helpers.HandleErr(w, err, global.StatusBadRequest, 400)
		maxPrice = 0

	}



	Items,first, last, errs := aih.itemService.ItemsByFilter(limit, offsetValue, searchField, categories, brand, types, sort, sortWay, minPrice, maxPrice)
	if len(errs) > 0 {
		helpers.HandleErr(w, errs, global.StatusNotFound, http.StatusNotFound)
		return
	}
	response := make(map[string]interface{})
	response["items"]= Items
	response["first"]=first
	response["last"]=last


	output, err := json.MarshalIndent(response, "", "\t\t")
	if err != nil {
		helpers.HandleErr(w, err, global.StatusInternalServerError, 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_,_ = w.Write(output)
	return

}

// GetItems handles GET /v1/admin/Items request
func (aih *AdminItemHandler) SearchItems(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	query := r.URL.Query()

	limit, err := strconv.Atoi(query.Get("limit"))
	if err != nil {
		helpers.HandleErr(w, err, global.StatusBadRequest, 400)
		return
	}
	if limit > global.MaxItemLimit {
		limit = global.MaxItemLimit
	}

	page, err := strconv.Atoi(query.Get("page"))
	if err != nil {
		helpers.HandleErr(w, err, global.StatusBadRequest, 400)
		page = global.DefaultPage
		return
	}

	var offset = page * limit

	Items, errs := aih.itemService.Items(limit, offset)
	if len(errs) > 0 {
		helpers.HandleErr(w, errs, global.StatusNotFound, http.StatusNotFound)
		return
	}

	output, err := json.MarshalIndent(Items, "", "\t\t")
	if err != nil {
		helpers.HandleErr(w, err, global.StatusInternalServerError, 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_,_ = w.Write(output)
	return

}

