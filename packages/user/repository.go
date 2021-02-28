package user

import (
	"context"
	"github.com/birukbelay/item/entity"
)

// UserRepository specifies application user related database operations
type UserRepository interface {
	GetUsers(ctx context.Context) ([]entity.User, []error)
	GetUser(ctx context.Context, id string) (*entity.User, []error)
	UpdateUser(ctx context.Context, user *entity.User) (*entity.User, []error)
	DeleteUser(ctx context.Context, id string) (*entity.User, []error)
	StoreUser(ctx context.Context, user *entity.User) (*entity.User, []error)

	UserByName(ctx context.Context, name string) (*entity.User, []error)
	UserByPhone(ctx context.Context, phone string) (*entity.User, []error)
	UserByEmail(ctx context.Context, email string) (*entity.User, []error)
	PhoneExists(ctx context.Context, phone string) bool
	EmailExists(ctx context.Context, email string) bool

}

