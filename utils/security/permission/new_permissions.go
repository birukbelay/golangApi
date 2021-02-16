package permission

import "strings"




type methods []string

type Permitted struct{
	roles map[string]methods
}

type paths map[string]Permitted

//type roles map[string][]string
//
//type paths map[string][]roles
//
//var role1= roles{"user":[]string{"post", "get"}}
//
//var p1 = paths{
//	"/items": []roles{
//		{ "client":[]string{"post","get"}},
//		{ "users":[]string{"post","get"}},
//		role1	},
//}

var allowed = paths{
	"/api/logout": Permitted{
		roles: map[string]methods{
			"CLIENT":{"POST","GET"},
			"client":{"POST"},
		},
	},
	"/api/items": Permitted{
		roles: map[string]methods{
			"ADMIN":{"POST", "GET"},
			"CLIENT":{ "GET"},
		},
	},
	"/api/items/create": Permitted{
		roles: map[string]methods{
			"CLIENT":{},
			"ADMIN":{"POST"},
		},
	},
}



func HavePermission(path string, role string, method string) bool {
	per := allowed[path]
	methodArr, ok:= per.roles[role]
	if !ok{
		return false
	}
	for _, m := range methodArr {
		if strings.ToUpper(m) == strings.ToUpper(method) {
			return true
		}
	}
	return false
}
