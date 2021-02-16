package jsonV

import (
	"encoding/json"
	"net/http"

	"github.com/birukbelay/item/utils/global"
	"github.com/birukbelay/item/utils/helpers"
	"github.com/birukbelay/item/utils/validators"
)

func JsonLoginValidator(w http.ResponseWriter, body []byte) bool {

	juser := JsonInput{}
	err2 := json.Unmarshal(body, &juser)
	if err2 != nil {
		helpers.HandleErr(w, err2, global.StatusBadRequest, 400)
		return false
	}

	newInput := JInput{Values: juser, VErrors: ValidationErrors{}}

	newInput.Required("info", "info_type", "password")

	if newInput.Values.Get("info_type").(string)=="email" {
		newInput.MatchesPattern("info", validators.EmailRX)
	}else if newInput.Values.Get("info_type").(string)=="phone" {
		newInput.MatchesPattern("info", validators.PhoneRX)
	}else{
		newInput.VErrors.Add("info","info_type not recognized")
	}

	newInput.MinLength("password", 6)

	if !newInput.Valid() {
		helpers.RenderResponse(w, newInput.VErrors, global.Validation, http.StatusBadRequest)
		return false
	}
	return true
}

func JsonSignupValidator(w http.ResponseWriter, body []byte) bool {

	juser := JsonInput{}
	err2 := json.Unmarshal(body, &juser)
	Eror2 := helpers.ErrSingle(w, err2, 404)
	if Eror2 != nil {
		return false
	}

	singupInput := JInput{Values: juser, VErrors: ValidationErrors{}}

	singupInput.Required("fullname", "email", "password", "confirmpassword")
	singupInput.MatchesPattern("email", validators.EmailRX)
	singupInput.MatchesPattern("phone", validators.PhoneRX)
	singupInput.MinLength("password", 8)
	singupInput.PasswordMatches("password", "confirmpassword")

	if !singupInput.Valid() {

		helpers.RenderResponse(w, singupInput.VErrors, global.Validation, http.StatusBadRequest)

		return false
	}
	return true
}

// Validate the form contents
//singnUpForm := form.Input{Values: r.PostForm, VErrors: form.ValidationErrors{}}
