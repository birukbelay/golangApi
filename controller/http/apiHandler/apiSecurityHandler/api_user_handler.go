package apiSecurityhandler

import (
	"encoding/json"
	"github.com/birukbelay/item/utils/global"
	"github.com/birukbelay/item/utils/validators/FormValidators"
	"github.com/julienschmidt/httprouter"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	user2 "github.com/birukbelay/item/packages/user"
	"github.com/birukbelay/item/utils/helpers"
	//my impors
	"github.com/birukbelay/item/entity"
	//rtoken2 "github.com/birukbelay/items/utils/security/rtoken"
	"github.com/birukbelay/item/utils/validators/form"
)


// UserHandler templateHandler handles user related requests
type UserHandler struct {
	userService user2.UserService
	signKey []byte
	loggedInUser   *entity.User
}

// NewUserHandler returns new UserHandler object
func NewUserHandler( usrServ user2.UserService, sKey []byte) *UserHandler {

	return &UserHandler{userService: usrServ,  signKey: sKey}
}



// AdminUsers handles Get /admin/users request
func (uh *UserHandler) AdminUsers(w http.ResponseWriter, r *http.Request) {

	users, errs := uh.userService.GetUsers()
	if len(errs) > 0 {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}

	output, err := json.MarshalIndent(users, "", "\t\t")

	Err := helpers.ErrSingle(w, err, 404)
	if Err != nil {
		return
	}

	w.Write(output)

}

// AdminUsersNew handles GET/POST /admin/users/new request creates new User With Role
func (uh *UserHandler) AdminUsersNew(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

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

		pExists := uh.userService.PhoneExists(r.FormValue("phone"))
		if pExists {
			accountForm.VErrors.Add("phone", "Phone Already Exists")
			helpers.RenderResponse(w, accountForm.VErrors, global.Validation, http.StatusBadRequest)
			return
		}
		eExists := uh.userService.EmailExists(r.FormValue("email"))
		if eExists {
			accountForm.VErrors.Add("email", "Email Already Exists")
			helpers.RenderResponse(w, accountForm.VErrors, global.Validation, http.StatusBadRequest)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), 12)
		if err != nil {
			accountForm.VErrors.Add("password", "Password Could not be stored")
			helpers.RenderResponse(w, accountForm.VErrors, global.Validation, http.StatusBadRequest)
			return
		}

	//TODO fetch the role from the Database
	roleID := r.FormValue("role")
	if (roleID == ""){
		roleID=entity.Client
	}


	var roles []string
	roles = append(roles, roleID)

		user := &entity.User{
			Username: r.FormValue("username"),
			Email:    r.FormValue("email"),
			Phone:    r.FormValue("phone"),
			Password: string(hashedPassword),
			Role:     roleID,
			Roles: roles,
		}
	//user.ID = primitive.NewObjectID()
		_, errs := uh.userService.StoreUser(user)
		if len(errs) > 0 {
			helpers.RenderResponse(w, errs, global.UserNotCreated, http.StatusInternalServerError)
			return
		}

	helpers.RenderResponse(w, user, global.UserCreated, http.StatusOK)
	return

}



// AdminUsersUpdate handles GET/POST /admin/users/update?id={id} request
func (uh *UserHandler) AdminUsersUpdate(w http.ResponseWriter, r *http.Request) {

		// Parse the form
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
		//helpers.LogValue("UserFormErrors", accountForm.VErrors)
		return
	}


	userID := r.FormValue("userid")
	phone:=r.FormValue("phone")
	email:=r.FormValue("email")
	username:= r.FormValue("username")
	roleID := r.FormValue("role")


	user, errs := uh.userService.GetUser(userID)
	if len(errs) > 0 {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}


	if user.Phone != phone {
		pExists := uh.userService.PhoneExists(r.FormValue("phone"))
		if pExists {
			accountForm.VErrors.Add("phone", "Phone Already Exists")
			//_ = uh.tmpl.ExecuteTemplate(w, "admin.user.update.layout", upAccForm)
			return
		}
	}
	if email!=user.Email {
		eExists := uh.userService.EmailExists(r.FormValue("email"))
		if eExists {
			accountForm.VErrors.Add("email", "Email Already Exists")
			helpers.RenderResponse(w, accountForm.VErrors, global.Validation, http.StatusBadRequest)
			return
		}
	}


		usr := &entity.User{
			ID:       user.ID,
			Username: username,
			Email:    email,
			Phone:    phone,
			Password: user.Password,
			Role:     roleID,
		}
		_, errs = uh.userService.UpdateUser(usr)
		if len(errs) > 0 {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}


}

// AdminUsersDelete handles Delete /admin/users/delete?id={id} request
func (uh *UserHandler) AdminUsersDelete(w http.ResponseWriter, r *http.Request) {
		idRaw := r.URL.Query().Get("id")
		_, errs := uh.userService.DeleteUser(idRaw)
		if len(errs) > 0 {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		}
}
