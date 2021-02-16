package service

import (
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
func (us *UserService) GetUsers() ([]entity.User, []error) {
	usrs, errs := us.userRepo.GetUsers()
	if len(errs) > 0 {
		return nil, errs
	}
	return usrs, errs
}

// GetUser retrieves an application user by its id
func (us *UserService) GetUser(id string) (*entity.User, []error) {
	usr, errs := us.userRepo.GetUser(id)
	if len(errs) > 0 {
		return nil, errs
	}
	return usr, errs
}
// RoleByName returns a user identified by its phone
func (us *UserService) UserByPhone(phone string) (*entity.User, []error) {
	user, errs := us.userRepo.UserByPhone(phone)
	if len(errs) > 0 {
		return nil, errs
	}
	return user, errs
}
// RoleByName returns a user identified by its phone
func (us *UserService) UserByName(phone string) (*entity.User, []error) {
	user, errs := us.userRepo.UserByName(phone)
	if len(errs) > 0 {
		return nil, errs
	}
	return user, errs
}


// UserByEmail retrieves an application user by its email address
func (us *UserService) UserByEmail(email string) (*entity.User, []error) {
	usr, errs := us.userRepo.UserByEmail(email)
	if len(errs) > 0 {
		return nil, errs
	}
	return usr, errs
}

// UpdateUser updates  a given application user
func (us *UserService) UpdateUser(user *entity.User) (*entity.User, []error) {
	usr, errs := us.userRepo.UpdateUser(user)
	if len(errs) > 0 {
		return nil, errs
	}
	return usr, errs
}

// DeleteUser deletes a given application user
func (us *UserService) DeleteUser(id string) (*entity.User, []error) {
	usr, errs := us.userRepo.DeleteUser(id)
	if len(errs) > 0 {
		return nil, errs
	}
	return usr, errs
}

// StoreUser stores a given application user
func (us *UserService) StoreUser(user *entity.User) (*entity.User, []error) {
	usr, errs := us.userRepo.StoreUser(user)
	if len(errs) > 0 {
		return nil, errs
	}
	return usr, errs
}

// PhoneExists check if there is a user with a given phone number
func (us *UserService) PhoneExists(phone string) bool {
	exists := us.userRepo.PhoneExists(phone)
	return exists
}

// EmailExists checks if there exist a user with a given email address
func (us *UserService) EmailExists(email string) bool {
	exists := us.userRepo.EmailExists(email)
	return exists
}
