package service

import (
	"ecomm/internal/domain"
	"ecomm/internal/dto"
	"ecomm/internal/repository"
	"errors"

	"ecomm/internal/helper"
)

type UserService struct {
	// Add necessary fields like repository, logger, etc.
	Repo repository.UserRepository
	Auth helper.Auth
}

// receiver function
// if used pointer for domain.User then need to return nil but if used without pointer then return empty interface/object

// in future we can add more fields like phone, name etc so it is better to use struct
// so instead of passing multiple parameters we can pass single input any
// any is alias for empty interface{}
func (us UserService) Signup(input dto.UserSignup) (string, error) {

	//hash the password before saving to db
	hashedPassword, err := us.Auth.CreateHashedPassword(input.Password)
	if err != nil {
		return "", err
	}

	user, err := us.Repo.CreateUser(domain.User{
		Email:    input.Email,
		Password: hashedPassword,
		Phone:    input.Phone,
	})
	return us.Auth.GenerateToken(user.ID, user.Email, user.UserType)
}

// find user by email
func (us UserService) FindUserByEmail(email string) (*domain.User, error) {
	// Implement the logic to find a user by email.
	//some business logic and database calls
	user, err := us.Repo.FindUser(email)
	return &user, err
}

func (us UserService) Login(email string, password string) (string, error) {
	// Implement the logic to sign up a user.
	//some business logic and database calls

	//call function find user by email
	user, err := us.FindUserByEmail(email)

	if err != nil {
		return "", errors.New("user not exist with this email")
	}

	//compare the password ,generate and return the token
	err = us.Auth.VerifyPassword(password, user.Password)

	if err != nil {
		return "", errors.New("invalid password")
	}
	//generate token and return
	return us.Auth.GenerateToken(user.ID, user.Email, user.UserType)

}

func (us UserService) GetVerificationCode(e domain.User) (int, error) {
	// Implement the logic to sign up a user.
	//some business logic and database calls
	return 0, nil
}
func (us UserService) VerifyCode(id uint, code int) error {
	// Implement the logic to sign up a user.
	//some business logic and database calls
	return nil
}
func (us UserService) CreateProfile(id uint, input any) error {

	return nil
}

// find the user by id and return the user profile but sometime
// we are using pointer bcoz at any point of time we need to edit the specific profile so taht why we are returning kind of pointer
func (us UserService) GetProfile(id uint) (*domain.User, error) {
	// Implement the logic to get a user profile.
	//some business logic and database calls
	return nil, nil
}
func (us UserService) UpdateProfile(id uint, input any) error {
	// Implement the logic to update a user profile.
	//some business logic and database calls
	return nil
}

func (us UserService) BecomeSeller(id uint, input any) (string, error) {

	return "", nil
}

// return bunch of card item so return slice of interface
func (us UserService) FindCart(id uint) ([]interface{}, error) {

	return nil, nil
}

// input as product properties and user info to update the cart
func (us UserService) CreateCart(input any, user domain.User) ([]interface{}, error) {

	return nil, nil
}

// just find the user whether the user have cart or not
func (us UserService) CreateOrder(user domain.User) (int, error) {

	return 0, nil
}

// accept the user id and find out the orders of that user
func (us UserService) GetOrders(user domain.User) ([]interface{}, error) {

	return nil, nil
}

// order id and user id
func (us UserService) GetOrderById(id uint, uId uint) (interface{}, error) {

	return nil, nil
}
