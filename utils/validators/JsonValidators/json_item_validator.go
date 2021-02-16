package jsonV

import (
	"encoding/json"
	"net/http"

	"github.com/birukbelay/item/utils/global"
	"github.com/birukbelay/item/utils/helpers"
)

func JsonItemValidator(w http.ResponseWriter,body []byte) bool  {

	jitems:=JsonInput{}
	err2 := json.Unmarshal(body, &jitems)
	if err2 != nil {
		helpers.HandleErr(w, err2, global.JsonUnmarshal, http.StatusInternalServerError)
		return false
	}

	newInput := JInput{Values: jitems, VErrors: ValidationErrors{}}

	newInput.Required("name", "price", "description", "image","categories")
	newInput.MinLength("description", 10)

	if !newInput.Valid() {

		helpers.RenderResponse(w, newInput.VErrors, global.Validation, http.StatusBadRequest)
		return false

	}
	return  true
}


func JsonCategoriesValidator(values JsonInput) (map [string][]string, bool)  {

	newCatForm := JInput{Values: values, VErrors: ValidationErrors{}}

	newCatForm.Required("catname", "catdesc")
	newCatForm.MinLength("catdesc", 10)

	if !newCatForm.Valid() {
		//ach.tmpl.ExecuteTemplate(w, "admin.categ.new.layout", newCatForm)
		return newCatForm.VErrors, false
	}
	return newCatForm.VErrors, true
}
