package user

 import "github.com/birukbelay/item/entity"

// UserService specifies application user related services
type UserService interface {
	GetUsers() ([]entity.User, []error)
	GetUser(id string) (*entity.User, []error)
	DeleteUser(id string) (*entity.User, []error)
	StoreUser(user *entity.User) (*entity.User, []error)

	UserByName(name string) (*entity.User, []error)
	UserByPhone(phone string) (*entity.User, []error)
	UserByEmail(email string) (*entity.User, []error)
	UpdateUser(user *entity.User) (*entity.User, []error)

	PhoneExists(phone string) bool
	EmailExists(email string) bool

}

