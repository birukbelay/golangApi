package productHandler


import (
	"encoding/json"
	"fmt"
	"github.com/birukbelay/item/entity"
	"github.com/birukbelay/item/utils/validators/FormValidators"
	"net/http"
	"net/url"
	"strconv"

	"github.com/julienschmidt/httprouter"
	//"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/birukbelay/item/models/items"
	"github.com/birukbelay/item/utils/global"
	"github.com/birukbelay/item/utils/helpers"
)

// AdminCategoriesHandler handles categories related http requests
type AdminCategoriesHandler struct {
	categoriesService items.CategoriesService
	//categoriesPerPage int
}

// NewAdminCategoriesHandler returns new AdminCategoriesHandler object
func NewAdminCategoriesHandler(genService items.CategoriesService) *AdminCategoriesHandler {
	return &AdminCategoriesHandler{categoriesService: genService}
}

// ==========================================  handlers ========================

// GetCategories handles GET /v1/admin/Categoriess request
func (aih *AdminCategoriesHandler) GetCategories(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	query := r.URL.Query()

	offsetValue :=query.Get("offsetValue")



	limit, err := strconv.Atoi(query.Get("limit"))
	if err != nil {
		helpers.HandleErr(w, err, global.StatusBadRequest, 400)
		limit = global.DefaultCategoriesLimit
	}
	if limit > global.MaxCategoriesLimit {
		limit = global.MaxCategoriesLimit
	}


	Categories, errs := aih.categoriesService.Categories(limit,offsetValue)
	if len(errs) > 0 {
		helpers.HandleErr(w, errs, global.StatusNotFound, http.StatusNotFound)
		return
	}

	output, err := json.MarshalIndent(Categories, "", "\t\t")
	if err != nil {
		helpers.HandleErr(w, err, global.StatusInternalServerError, 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
	return

}

// GetSingleCategories handles GET /categories/create:id request
func (aih *AdminCategoriesHandler) GetSingleCategories(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {

	id := ps.ByName("id")
	fmt.Println(id)

	// calling the service
	categories, errs := aih.categoriesService.Category(id)
	if len(errs) > 0 {
		helpers.HandleErr(w, errs, global.StatusNotFound, http.StatusNotFound)
		return
	}

	output, err := json.MarshalIndent(categories, "", "\t\t")

	if err != nil {
		helpers.HandleErr(w, err, global.StatusBadRequest, 400)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
	return
}

// CreateCategories handles POST /v1/admin/Categories request
func (aih *AdminCategoriesHandler) CreateCategories(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	if err := r.ParseMultipartForm(global.MaxUploadSize); err != nil {
		//fmt.Printf("Could not parse multipart form: %v\n", err)
		helpers.RenderResponse(w, err, global.ParseFile, http.StatusBadRequest)
		return
	}

	// form is found in github.com/birukbelay/items/utils/validators/form
	VErrors, valid := FormValidators.FormCategoriesValidator(r.PostForm)
	if !valid {
		helpers.RenderResponse(w, VErrors, global.ParseFile, http.StatusBadRequest)
		return
	}
	categories, er := InitiateCategories(r.PostForm)
	if er!=nil{
		helpers.RenderResponse(w, er, global.CategoriesInitialization, http.StatusBadRequest)
	}
	//TODO fill the rest of the categories fields


	//categories.ID = primitive.NewObjectID()

	// IMAGE UPLOAD
	img, err, status, statusCode := helpers.UploadFile(r, false, "", "Categoriess" )
	if err != nil {
		helpers.RenderResponse(w, err, status, statusCode)
		return
	}
	categories.Image = img

	// calling the service
	gen, errs := aih.categoriesService.StoreCategories(categories)
	if len(errs) > 0 {
		helpers.HandleErr(w, errs, global.StatusInternalServerError, 404)
		return
	}

	helpers.RenderResponse(w, gen, global.Success, http.StatusCreated)

	//p := fmt.Sprintf("/v1/admin/Categories/%d", categories.ID)
	//w.Header().Set("Location", p)
	return
}

// UpdateCategories handles PUT /categories/update/:id request
func (aih *AdminCategoriesHandler) UpdateCategories(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	id := ps.ByName("id")
	//query := r.URL.Query()

	//hasImage, err := strconv.ParseBool(query.Get("hasImage"))
	//if err != nil {
	//	helpers.HandleErr(w, err, global.StatusBadRequest, 400)
	//	return
	//}

	gen, errs := aih.categoriesService.Category(id)
	if len(errs) > 0 {
		helpers.HandleErr(w, errs, global.StatusNotFound, http.StatusNotFound)
		return
	}


	if err := r.ParseMultipartForm(global.MaxUploadSize); err != nil {
		//fmt.Printf("Could not parse multipart form: %v\n", err)
		helpers.RenderResponse(w, err, global.ParseFile, http.StatusBadRequest)
		return
	}

	// form is found in github.com/birukbelay/items/utils/validators/form
	VErrors, valid := FormValidators.FormCategoriesValidator(r.PostForm)
	if !valid {
		helpers.RenderResponse(w, VErrors, global.ParseFile, http.StatusBadRequest)
		return
	}

	categories, er := InitiateCategories(r.PostForm)
	if er!=nil{
		helpers.RenderResponse(w, er, global.CategoriesInitialization, http.StatusBadRequest)
	}

	imageChanged , er := strconv.ParseBool(r.PostForm.Get("imageChanged"))
	if er!=nil{
		imageChanged =false
	}
	userSess, _ := r.Context().Value(entity.CtxUserSessionKey).(*entity.User)


	image:= gen.Image
	if imageChanged {
		img, err, status, statusCode := helpers.UploadFile(r,true, image, userSess.ID.Hex())
		if err != nil {
			helpers.RenderResponse(w, err, status, statusCode)
			img = gen.Image
			return
		}
		categories.Image = img
	//TODO make a function to Change the image, Delete the Image
	}else {
		categories.Image=  gen.Image
	}


	//oid, _ := primitive.ObjectIDFromHex(id)
	//categories.ID =  oid
	//TODo make Updated count and UserId



	// calling the service
	categories, errs = aih.categoriesService.UpdateCategories(categories)
	if len(errs) > 0 {
		helpers.HandleErr(w, errs, global.StatusNotFound, 404)
		return
	}

	helpers.RenderResponse(w, categories, global.Success, http.StatusOK)
	return
}

// DeleteCategories handles DELETE /categories/remove/:id request
func (aih *AdminCategoriesHandler) DeleteCategories(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	id := ps.ByName("id")

	// calling the service
	categories, errs := aih.categoriesService.Category(id)
	//TODO do stg with categories
	fmt.Println(categories)

	if len(errs) > 0 {
		helpers.HandleErr(w, errs, global.StatusNotFound, http.StatusNotFound)
		return
	}
	// calling the service
	_, errs = aih.categoriesService.DeleteCategories(id)
	if len(errs) > 0 {
		helpers.HandleErr(w, errs, global.StatusInternalServerError, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	_,_ = w.Write([]byte("Deleted"))
	return
}

func InitiateCategories(values url.Values) (*entity.Categories, error) {

	categories := &entity.Categories{
		Name:         values.Get("name"),
		Description:  values.Get("description"),
	}
	return categories, nil
}
