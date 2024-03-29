package productHandler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/birukbelay/item/entity"
	"github.com/birukbelay/item/utils/validators/FormValidators"

	"github.com/birukbelay/item/packages/items"
	"github.com/birukbelay/item/utils/global"
	"github.com/birukbelay/item/utils/helpers"
	"github.com/julienschmidt/httprouter"
)

// AdminItemHandler handles items related http requests
type AdminItemHandler struct {
	itemService items.ItemService
	//itemPerPage int
}

// NewAdminItemHandler returns new AdminItemHandler object
func NewAdminItemHandler(itmService items.ItemService) *AdminItemHandler {
	return &AdminItemHandler{itemService: itmService}
}

// ==========================================  handlers ========================

// GetItems handles GET /v1/admin/Items request
func (aih *AdminItemHandler) GetItems(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	contxt := r.Context()
	var ctx, _ = context.WithTimeout(contxt, 30*time.Second)

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

	Items, errs := aih.itemService.Items(ctx, limit, offset)
	if len(errs) > 0 {
		helpers.LogTrace("item fetch err", err)
		helpers.HandleErr(w, errs, global.StatusNotFound, http.StatusNotFound)
		return
	}

	output, err := json.MarshalIndent(Items, "", "\t\t")
	if err != nil {
		helpers.HandleErr(w, err, global.StatusInternalServerError, 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
	return

}

// GetSingleItem handles GET /items/create:id request
func (aih *AdminItemHandler) GetSingleItem(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	id := ps.ByName("id")
	fmt.Println(id)
	contxt := r.Context()
	var ctx, _ = context.WithTimeout(contxt, 30*time.Second)

	// calling the service
	item, errs := aih.itemService.Item(ctx, id)
	if len(errs) > 0 {
		helpers.RenderResponse(w, errs, global.StatusNotFound, http.StatusNotFound)
		return
	}

	output, err := json.MarshalIndent(item, "", "\t\t")

	if err != nil {
		helpers.HandleErr(w, err, global.StatusBadRequest, 400)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
	return
}

// InitiateItem handles POST /v1/admin/Items request
func (aih *AdminItemHandler) CreateItem(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	contxt := r.Context()
	var ctx, _ = context.WithTimeout(contxt, 30*time.Second)

	if err := r.ParseMultipartForm(global.MaxUploadSize); err != nil {
		//fmt.Printf("Could not parse multipart form: %v\n", err)
		helpers.RenderResponse(w, err, global.ParseFile, http.StatusBadRequest)
		return
	}

	// form is found in github.com/birukbelay/items/utils/validators/form
	valid, VErrors := FormValidators.FormItemValidator(r.PostForm)
	if !valid {
		helpers.RenderResponse(w, VErrors, global.Validation, http.StatusBadRequest)
		return
	}

	item, err := InitiateItem(r.PostForm)
	if err != nil {
		helpers.RenderResponse(w, err, global.ItemInitialization, http.StatusBadRequest)
		return
	}

	//items.ID = primitive.NewObjectID()

	// IMAGE UPLOAD
	fmt.Println("beSess..")
	//userSess, _ := r.Context().Value(entity.CtxUserSessionKey).(*entity.User)

	// IMAGE UPLOAD
	// fmt.Println("beUpl..",userSess)

	img, er, status, statusCode := helpers.UploadFile(r, false, "", "items")
	if er != nil {
		helpers.RenderResponse(w, err, status, statusCode)
		return
	}
	item.Image = img

	// calling the service
	item, errs := aih.itemService.StoreItem(ctx, item)
	if len(errs) > 0 {
		helpers.HandleErr(w, errs, global.StatusInternalServerError, 404)
		return
	}

	helpers.RenderResponse(w, item, global.Success, http.StatusCreated)

	//p := fmt.Sprintf("/v1/admin/Items/%d", items.ID)
	//w.Header().Set("Location", p)
	return
}

// UpdateItem handles PUT /items/update/:id request
func (aih *AdminItemHandler) UpdateItem(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	contxt := r.Context()
	var ctx, _ = context.WithTimeout(contxt, 30*time.Second)

	id := ps.ByName("id")
	helpers.LogTrace("updateId", id)
	//query := r.URL.Query()

	//hasImage, err := strconv.ParseBool(query.Get("hasImage"))
	//if err != nil {
	//	helpers.HandleErr(w, err, global.StatusBadRequest, 400)
	//	return
	//}

	if err := r.ParseMultipartForm(global.MaxUploadSize); err != nil {
		//fmt.Printf("Could not parse multipart form: %v\n", err)
		helpers.RenderResponse(w, err, global.ParseFile, http.StatusBadRequest)
		return
	}
	valid, VErrors := FormValidators.FormItemValidator(r.PostForm)
	if !valid {
		helpers.RenderResponse(w, VErrors, global.Validation, http.StatusBadRequest)
		return
	}

	itm, errs := aih.itemService.Item(ctx, id)
	if len(errs) > 0 {
		helpers.HandleErr(w, errs, global.StatusNotFound, http.StatusNotFound)
		return
	}

	item, err := InitiateItem(r.PostForm)
	if err != nil {
		helpers.RenderResponse(w, err, global.ItemInitialization, http.StatusBadRequest)
		return
	}
	itm.Price = item.Price
	itm.Categories = item.Categories
	itm.Name = item.Name
	itm.Description = item.Description

	helpers.LogTrace("foundItem FOr Update", itm)

	imageChanged, er := strconv.ParseBool(r.PostForm.Get("imageChanged"))
	if er != nil {
		imageChanged = false
	}
	//userSess, _ := r.Context().Value(entity.CtxUserSessionKey).(*entity.User)

	image := itm.Image

	if imageChanged {
		img, err, status, statusCode := helpers.UploadFile(r, true, image, "items")
		if err != nil {
			helpers.LogTrace("UpdateImgErr", err)
			helpers.RenderResponse(w, err, status, statusCode)
			img = itm.Image
			return
		}
		itm.Image = img
		//TODO make a function to Change the image, Delete the Image
	} else {
		itm.Image = itm.Image
	}

	// calling the service
	item, errs = aih.itemService.UpdateItem(ctx, itm)
	if len(errs) > 0 {
		helpers.HandleErr(w, errs, global.StatusNotFound, 404)
		return
	}

	helpers.RenderResponse(w, itm, global.Success, http.StatusOK)
	return
}

// DeleteItem handles DELETE /items/remove/:id request
func (aih *AdminItemHandler) DeleteItem(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	contxt := r.Context()
	var ctx, _ = context.WithTimeout(contxt, 30*time.Second)

	id := ps.ByName("id")

	// calling the service
	item, errs := aih.itemService.Item(ctx, id)
	fmt.Println(item)

	if len(errs) > 0 {
		helpers.RenderResponse(w, errs, global.StatusNotFound, http.StatusNotFound)
		return
	}
	// calling the service
	count, errs := aih.itemService.DeleteItem(ctx, id)
	if len(errs) > 0 {
		helpers.HandleErr(w, errs, global.StatusInternalServerError, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte(string(count)))
	return
}

func InitiateItem(values url.Values) (*entity.Item, error) {

	year, err := strconv.ParseInt(values.Get("price"), 10, 32)
	if err != nil {
		return nil, err
	}

	category := values.Get("category")
	var categories []string
	categories = append(categories, category)
	item := &entity.Item{
		Name:        values.Get("name"),
		Description: values.Get("description"),
		//Image:        values.Get("image"),
		Categories: categories,
		Type:       values.Get("type"),

		Price: int(year),
	}
	return item, nil
}
