package apiSecurityhandler

import (
	"context"
	"github.com/birukbelay/item/entity"
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/birukbelay/item/utils/helpers"
	permission2 "github.com/birukbelay/item/utils/security/permission"
	"github.com/birukbelay/item/utils/security/rtoken"
)


func (uh *UserHandler) loggedIn(r *http.Request) (*rtoken.CustomClaims, bool) {

	if r.Header["Token"] == nil {
		return nil, false
	}

	signKey := uh.signKey
	if r.Header["Token"] != nil{
		token:= r.Header["Token"][0]
		claims, ok, err := rtoken.Valid(token, signKey)
		if !ok || (err != nil) {
			helpers.LogValue("TokenError",err)
			return nil, false
		}
		//user := &claims.User
		//uh.loggedInUser=user
		return claims, true
	}
	return nil, false

}


// Authenticated checks if a user is authenticated to access a given route
func (uh *UserHandler) Authenticated(next httprouter.Handle) httprouter.Handle {
	fn := func(w http.ResponseWriter, r *http.Request,  ps httprouter.Params) {
	claims,	ok := uh.loggedIn(r)
		if !ok {
			helpers.LogValue("Authenticated function","")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			//http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		uh.loggedInUser= &claims.User


		ctx := context.WithValue(r.Context(), entity.CtxUserSessionKey,  claims.User)
		next(w, r.WithContext(ctx), ps)
	}
	return fn
}

// Authorized checks if a user has proper authority to access a give route
func (uh *UserHandler) Authorized(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {


		if uh.loggedInUser == nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}


		userSess, _ := r.Context().Value(entity.CtxUserSessionKey).(*entity.User)

		roles := userSess.Roles
		for _, role := range roles {

			permitted := permission2.HavePermission(r.URL.Path, role, r.Method)
			if permitted {
				next(w, r, ps)
				return
			}
		}

		helpers.LogValue("roles Not Permited ...  ${role}","")
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}
}


