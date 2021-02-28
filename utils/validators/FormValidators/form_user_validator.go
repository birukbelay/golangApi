package FormValidators

import (
	"github.com/birukbelay/item/utils/validators"
	"github.com/birukbelay/item/utils/validators/form"
)

func FormUserValidator(accountForm *form.Input)  bool  {

	// Validate the form contents
	//accountForm := Input{Values: values, VErrors: ValidationErrors{}}

	accountForm.Required("username", "email", "password", "confirmpassword")
	accountForm.MatchesPattern("email", validators.EmailRX)
	accountForm.MatchesPattern("phone", validators.PhoneRX)
	accountForm.MinLength("password", 6)
	accountForm.PasswordMatches("password", "confirmpassword")
	//accountForm.CSRF = token
	// If there are any errors, redisplay the signup form.
	if !accountForm.Valid() {
		return  false
	}
	//fmt.Println("valid")
	return  true
}
func FormUserUpdateValidator(accountForm *form.Input)  bool  {

	// Validate the form contents
	//accountForm := Input{Values: values, VErrors: ValidationErrors{}}

	accountForm.Required("username", "email",)
	accountForm.MatchesPattern("email", validators.EmailRX)
	accountForm.MatchesPattern("phone", validators.PhoneRX)
	accountForm.MinLength("password", 6)
	accountForm.PasswordMatches("password", "confirmpassword")
	//accountForm.CSRF = token
	// If there are any errors, redisplay the signup form.
	if !accountForm.Valid() {
		return  false
	}
	//fmt.Println("valid")
	return  true
}

func LoginValidator(accountForm *form.Input)  bool  {


	accountForm.Required("info", "info_type", "password")
	if accountForm.Values.Get("info_type")=="email" {
		accountForm.MatchesPattern("info", validators.EmailRX)
	}else if accountForm.Values.Get("info_type")=="phone" {
		accountForm.MatchesPattern("info", validators.PhoneRX)
	}else{
		accountForm.VErrors.Add("info","info_type not recognized")
	}
	accountForm.MinLength("password", 6)

	if !accountForm.Valid() {
		return  false
	}
	//fmt.Println("valid")
	return  true
}