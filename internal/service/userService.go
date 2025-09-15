package service

import (
	"ecomm/internal/domain"
	"ecomm/internal/dto"
	"log"
)

type UserService struct {
	// Add necessary fields like repository, logger, etc.

}

// receiver function
// if used pointer for domain.User then need to return nil but if used without pointer then return empty interface/object

// in future we can add more fields like phone, name etc so it is better to use struct
// so instead of passing multiple parameters we can pass single input any
// any is alias for empty interface{}
func (us UserService) Signup(input dto.UserSignup) (string, error) {
	// Implement the logic to sign up a user.
	//some business logic and database calls
	log.Println(input)
	return "this is my token as of now", nil
}

func (us UserService) Login(string, error) (string, error) {
	// Implement the logic to sign up a user.
	//some business logic and database calls
	return "", nil
}

func (us UserService) FindUserByEmail(email string) (*domain.User, error) {
	// Implement the logic to find a user by email.
	//some business logic and database calls
	return nil, nil
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
