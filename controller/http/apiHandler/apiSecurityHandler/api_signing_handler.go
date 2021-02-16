package apiSecurityhandler

import (
	"fmt"
	"github.com/birukbelay/item/utils/validators/FormValidators"
	"github.com/birukbelay/item/utils/validators/form"
	"net/http"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"

	//my imports
	"github.com/birukbelay/item/entity"
	"github.com/birukbelay/item/utils/global"
	"github.com/birukbelay/item/utils/helpers"
	rtoken2 "github.com/birukbelay/item/utils/security/rtoken"
)

func (uh *UserHandler) checkAdmin(rs []string) bool {
	for _, r := range rs {
		if strings.ToUpper(r) == strings.ToUpper("Admin") {
			return true
		}
	}
	return false
}


// ApiLogin  the POST /login requests
func (uh *UserHandler) ApiLogin(w http.ResponseWriter, r *http.Request,  _ httprouter.Params) {

	if err := r.ParseMultipartForm(global.MaxUploadSize); err != nil {
		//fmt.Printf("Could not parse multipart form: %v\n", err)
		helpers.RenderResponse(w,err, global.ParseFile, http.StatusBadRequest)
		return
	}

	// form is found in github.com/birukbelay/items/utils/validators/form
	accountForm := form.Input{Values: r.PostForm, VErrors: form.ValidationErrors{}}
	valid:= FormValidators.LoginValidator(&accountForm)
	if !valid{
		helpers.RenderResponse(w, accountForm.VErrors, global.Validation, http.StatusBadRequest)
		helpers.LogValue("UserFormErrors", accountForm.VErrors)
		return
	}

	loginData := &entity.LoginData{
		LoginInfo: r.FormValue("info"),
		InfoType: r.FormValue("info_type"),
		Password: r.FormValue("password"),
	}



	user := &entity.User{}

	if loginData.InfoType=="email"{
		usr, errs := uh.userService.UserByEmail(loginData.LoginInfo)
		if len(errs) > 0 {
			helpers.HandleErr(w, errs, global.EmailOrPassword, 400)
			return
		}
		user=usr
	}else if  loginData.InfoType=="phone"{
		usr, errs := uh.userService.UserByPhone(loginData.LoginInfo)
		if len(errs) > 0 {
			helpers.HandleErr(w, errs, global.EmailOrPassword, 400)
			return
		}
		user=usr
	}else{
		helpers.HandleErr(w, global.InvalidEmailOrPhone, global.InvalidData, 400)
	}



	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(r.FormValue("password")))
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			helpers.HandleErr(w, err, global.EmailOrPassword, 400)
			return
		}
		helpers.HandleErr(w, err, global.EmailOrPassword, 400)
		return
	}
	session, er := CreateSession(user)
	if er!=nil{
		helpers.RenderResponse(w, global.UserNotCreated, global.StatusInternalServerError, 500)
		return
	}
	user.Session = append(user.Session, *session)

	_, errs := uh.userService.UpdateUser(user)
	if len(errs) > 0 {
		helpers.RenderResponse(w, global.StatusInternalServerError, global.StatusInternalServerError, 500)
		return
	}

	uh.Aftermath(w ,user, session)

	helpers.RenderResponse(w, user, global.Success, http.StatusOK)

}



// Signup hanldes the GET/POST /signup requests
func (uh *UserHandler) ApiSignup(w http.ResponseWriter, r *http.Request,  _ httprouter.Params) {

	if err := r.ParseMultipartForm(global.MaxUploadSize); err != nil {
		//fmt.Printf("Could not parse multipart form: %v\n", err)
		helpers.RenderResponse(w,err, global.ParseFile, http.StatusBadRequest)
		return
	}

	// form is found in github.com/birukbelay/items/utils/validators/form
	accountForm := form.Input{Values: r.PostForm, VErrors: form.ValidationErrors{}}
	valid:= FormValidators.FormUserValidator(&accountForm)
	if !valid{
		helpers.RenderResponse(w, accountForm.VErrors, global.Validation, http.StatusBadRequest)
		helpers.LogValue("UserFormErrors", accountForm.VErrors)
		return
	}



	user := &entity.User{
		Username: r.FormValue("username"),
		Email:    r.FormValue("email"),
		Phone:    r.FormValue("phone"),
	}

	pExists := uh.userService.PhoneExists(user.Phone)
	if pExists {
		helpers.HandleErr(w, global.PhoneExists, global.PhoneExists, http.StatusBadRequest)
		return
	}
	eExists := uh.userService.EmailExists(user.Email)
	if eExists {
		helpers.HandleErr(w, global.EmailExists, global.EmailExists, http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), 12)
	if err != nil {
		helpers.HandleErr(w, global.Password, global.Password, http.StatusBadRequest)
		return
	}

	//role, errs := uh.userRole.RoleByName("USER")
	//if len(errs) > 0 {
	//	helpers.RenderResponse(w, global.Role, global.Role, http.StatusInternalServerError)
	//	return
	//}

	user.Password = string(hashedPassword)


	user.Roles = append(user.Roles, entity.Admin)
	user.Role = entity.Admin

	session, er := CreateSession(user)
	if er!=nil{
		helpers.RenderResponse(w, global.UserNotCreated, global.StatusInternalServerError, 500)
		return
	}
	user.Session = append(user.Session, *session)
//fmt.Println("...before Storing")

	usr, errs := uh.userService.StoreUser(user)
	if len(errs) > 0 {
		fmt.Println("...error")
		helpers.RenderResponse(w, global.StatusInternalServerError, global.StatusInternalServerError, 500)
		return
	}

	uh.Aftermath(w , user, session)

	helpers.RenderResponse(w, usr, global.UserCreated, http.StatusCreated)


}


func CreateSession(user *entity.User) (*entity.Session, error) {

	sessionUUID := rtoken2.CategoryteRandomID(32)
	signingString, err := rtoken2.CategoryteRandomString(32)
	if err != nil {
		return nil, err
	}

	signingKey := []byte(signingString)

	Sess:=&entity.Session{ UUID: sessionUUID, SigningKey: signingKey, LoginDate: time.Now()}
	return Sess, nil
}

func (uh *UserHandler)Aftermath(w http.ResponseWriter, user *entity.User, session * entity.Session){

	uh.loggedInUser = user


	tokenExpires := time.Now().Add(time.Hour * 24*15).Unix()

	claims := rtoken2.Claims(user, string(session.UUID), user.ID, tokenExpires)

	rtoken2.CreateToken(w, uh.signKey, claims)

}

