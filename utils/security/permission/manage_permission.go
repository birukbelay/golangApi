package permission

import "strings"
//==================================================================================================
//=============================================- this is the old Permission ---===================

type permission struct {
	roles   []string
	methods []string
}

// authority is a map of 'string' and the 'permission' struct
type authority map[string]permission

//............................................................................................
//............................................................................................





func checkRole(role string, roles []string) bool {
	for _, r := range roles {
		if strings.ToUpper(r) == strings.ToUpper(role) {
			return true
		}
	}
	return false
}

func checkMethod(method string, methods []string) bool {
	for _, m := range methods {
		if strings.ToUpper(m) == strings.ToUpper(method) {
			return true
		}
	}
	return false
}


// HasPermission checks if a given role has permission to access a given route for a given method
func HasPermission(path string, role string, method string) bool {
	///admin/categs/new
	if strings.HasPrefix(path, "/admin") {
		path = "/admin"
	}
	///order/1
	if strings.HasPrefix(path, "/order"){
		path = "/order"
	}
	// perm - is value of the key [path] -> this is getting the value from the 'authorites' map  by using the [path] key
	perm := authorities[path]

	checkedRole := checkRole(role, perm.roles)
	checkedMethod := checkMethod(method, perm.methods)
	if !checkedRole || !checkedMethod {
		return false
	}
	return true
}




// authorities is instance of the authority map
var authorities = authority{

	"/contact": permission{
		roles:   []string{"USER"},
		methods: []string{"GET", "POST"},
	},

	"/login": permission{
		roles:   []string{"USER"},
		methods: []string{"GET", "POST"},
	},
	"/about": permission{
		roles:   []string{"USER"},
		methods: []string{"GET"},
	},
	"/logout": permission{
		roles:   []string{"USER"},
		methods: []string{"POST"},
	},
	"/signup": permission{
		roles:   []string{"USER"},
		methods: []string{"GET", "POST"},
	},
	"/menu": permission{
		roles:   []string{"USER"},
		methods: []string{"GET"},
	},
	"/order": permission{
		roles:   []string{"USER"},
		methods: []string{"GET", "POST"},
	},
	"/admin": permission{
		roles:   []string{"ADMIN"},
		methods: []string{"GET", "POST"},
	},
}
