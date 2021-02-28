package apiSecurityhandler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/birukbelay/item/utils/global"
	"github.com/birukbelay/item/utils/validators/FormValidators"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"time"

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
func (uh *UserHandler) AdminUsers(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	contxt := r.Context()
	var ctx, _ = context.WithTimeout(contxt, 30*time.Second)
	users, errs := uh.userService.GetUsers(ctx)
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

func (uh *UserHandler) User(w http.ResponseWriter, r *http.Request,ps httprouter.Params){
	id := ps.ByName("id")
	fmt.Println(id)
	contxt := r.Context()
	var ctx, _ = context.WithTimeout(contxt, 30*time.Second)

	// calling the service
	item, errs := uh.userService.GetUser(ctx, id)
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


// AdminUsersNew handles GET/POST /admin/users/new request creates new User With Role
func (uh *UserHandler) AdminUsersNew(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	contxt := r.Context()
	var ctx, _ = context.WithTimeout(contxt, 30*time.Second)
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

		pExists := uh.userService.PhoneExists(ctx, r.FormValue("phone"))
		if pExists {
			accountForm.VErrors.Add("phone", "Phone Already Exists")
			helpers.RenderResponse(w, accountForm.VErrors, global.Validation, http.StatusBadRequest)
			return
		}
		eExists := uh.userService.EmailExists(ctx, r.FormValue("email"))
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
		_, errs := uh.userService.StoreUser(ctx, user)
		if len(errs) > 0 {
			helpers.RenderResponse(w, errs, global.UserNotCreated, http.StatusInternalServerError)
			return
		}

	helpers.RenderResponse(w, user, global.UserCreated, http.StatusOK)
	return

}



// AdminUsersUpdate handles GET/POST /admin/users/update?id={id} request
func (uh *UserHandler) AdminUsersUpdate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	contxt := r.Context()
	var ctx, _ = context.WithTimeout(contxt, 30*time.Second)
	id := ps.ByName("id")
		// Parse the form
		helpers.LogTrace("update", "here")
	if err := r.ParseMultipartForm(global.MaxUploadSize); err != nil {
		//fmt.Printf("Could not parse multipart form: %v\n", err)
		helpers.RenderResponse(w,err, global.ParseFile, http.StatusBadRequest)
		return
	}

	// form is found in github.com/birukbelay/items/utils/validators/form
	accountForm := form.Input{Values: r.PostForm, VErrors: form.ValidationErrors{}}
	valid:= FormValidators.FormUserUpdateValidator(&accountForm)
	if !valid{
		helpers.RenderResponse(w, accountForm.VErrors, global.Validation, http.StatusBadRequest)
		//helpers.LogValue("UserFormErrors", accountForm.VErrors)
		return
	}




	phone:=r.FormValue("phone")
	email:=r.FormValue("email")
	username:= r.FormValue("username")
	roleID := r.FormValue("role")


	user, errs := uh.userService.GetUser(ctx, id)
	if len(errs) > 0 {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}


	if user.Phone != phone {
		pExists := uh.userService.PhoneExists(ctx, r.FormValue("phone"))
		if pExists {
			accountForm.VErrors.Add("phone", "Phone Already Exists")
			//_ = uh.tmpl.ExecuteTemplate(w, "admin.user.update.layout", upAccForm)
			return
		}
	}
	if email!=user.Email {
		eExists := uh.userService.EmailExists(ctx, r.FormValue("email"))
		if eExists {
			accountForm.VErrors.Add("email", "Email Already Exists")
			helpers.RenderResponse(w, accountForm.VErrors, global.Validation, http.StatusBadRequest)
			return
		}
	}


		if username !=""{
			user.Username=username
		}
		if email !=""{
			user.Email=email
		}
		if phone !=""{
			user.Phone=phone
		}


		if roleID !=""{
			user.Role=roleID
		}


		user, errs = uh.userService.UpdateUser(ctx, user)
		if len(errs) > 0 {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	helpers.RenderResponse(w, user, global.UserCreated, http.StatusAccepted)
	return


}

// AdminUsersDelete handles Delete /admin/users/delete?id={id} request
func (uh *UserHandler) AdminUsersDelete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	contxt := r.Context()
	var ctx, _ = context.WithTimeout(contxt, 30*time.Second)
	id := ps.ByName("id")
		user, errs := uh.userService.DeleteUser(ctx, id)
		if len(errs) > 0 {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		}
	helpers.RenderResponse(w, user, global.Success, http.StatusNoContent)
	return
}
