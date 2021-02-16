package FormValidators

import (
	"github.com/birukbelay/item/utils/validators/form"
	"net/url"
)

func FormCategoriesValidator(values url.Values) (form.ValidationErrors, bool)  {

	newCategories := form.Input{Values: values, VErrors: form.ValidationErrors{}}
	newCategories.Required("name", "description")
	newCategories.MinLength("description", 10)

	// If there are any errors, redisplay the signup form.
	if !newCategories.Valid() {
		//ach.tmpl.ExecuteTemplate(w, "admin.categ.new.layout", newCategories)
		return newCategories.VErrors, false
	}
	return newCategories.VErrors, true
}

func FormItemValidator(values url.Values) (bool, form.ValidationErrors)  {

	newItemForm := form.Input{Values: values, VErrors: form.ValidationErrors{}}
	newItemForm.Required("name", "description","category", )
	newItemForm.MinLength("description", 10)


	if !newItemForm.Valid() {
		return false, newItemForm.VErrors
	}

	return true, newItemForm.VErrors
}


