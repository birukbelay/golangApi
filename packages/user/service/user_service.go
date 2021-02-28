package service

import (
	"context"
	user2 "github.com/birukbelay/item/packages/user"
	"github.com/birukbelay/item/entity"
)

// UserService implements user.UserService interface
type UserService struct {
	userRepo user2.UserRepository
}

// NewUserService  returns a new UserService object
func NewUserService(userRepository user2.UserRepository) user2.UserService {
	return &UserService{userRepo: userRepository}
}

// GetUsers returns all stored application users
func (us *UserService) GetUsers(ctx context.Context) ([]entity.User, []error) {
	usrs, errs := us.userRepo.GetUsers(ctx)
	if len(errs) > 0 {
		return nil, errs
	}
	return usrs, errs
}

// GetUser retrieves an application user by its id
func (us *UserService) GetUser(ctx context.Context, id string) (*entity.User, []error) {
	usr, errs := us.userRepo.GetUser(ctx, id)
	if len(errs) > 0 {
		return nil, errs
	}
	return usr, errs
}
// RoleByName returns a user identified by its phone
func (us *UserService) UserByPhone(ctx context.Context, phone string) (*entity.User, []error) {
	user, errs := us.userRepo.UserByPhone(ctx, phone)
	if len(errs) > 0 {
		return nil, errs
	}
	return user, errs
}
// RoleByName returns a user identified by its phone
func (us *UserService) UserByName(ctx context.Context, phone string) (*entity.User, []error) {
	user, errs := us.userRepo.UserByName(ctx, phone)
	if len(errs) > 0 {
		return nil, errs
	}
	return user, errs
}


// UserByEmail retrieves an application user by its email address
func (us *UserService) UserByEmail(ctx context.Context, email string) (*entity.User, []error) {
	usr, errs := us.userRepo.UserByEmail(ctx, email)
	if len(errs) > 0 {
		return nil, errs
	}
	return usr, errs
}

// UpdateUser updates  a given application user
func (us *UserService) UpdateUser(ctx context.Context, user *entity.User) (*entity.User, []error) {
	usr, errs := us.userRepo.UpdateUser(ctx, user)
	if len(errs) > 0 {
		return nil, errs
	}
	return usr, errs
}

// DeleteUser deletes a given application user
func (us *UserService) DeleteUser(ctx context.Context, id string) (*entity.User, []error) {
	usr, errs := us.userRepo.DeleteUser(ctx, id)
	if len(errs) > 0 {
		return nil, errs
	}
	return usr, errs
}

// StoreUser stores a given application user
func (us *UserService) StoreUser(ctx context.Context, user *entity.User) (*entity.User, []error) {
	usr, errs := us.userRepo.StoreUser(ctx, user)
	if len(errs) > 0 {
		return nil, errs
	}
	return usr, errs
}

// PhoneExists check if there is a user with a given phone number
func (us *UserService) PhoneExists(ctx context.Context, phone string) bool {
	exists := us.userRepo.PhoneExists(ctx, phone)
	return exists
}

// EmailExists checks if there exist a user with a given email address
func (us *UserService) EmailExists(ctx context.Context, email string) bool {
	exists := us.userRepo.EmailExists(ctx, email)
	return exists
}
