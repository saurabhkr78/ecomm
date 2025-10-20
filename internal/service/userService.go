package service

import (
	"ecomm/configs"
	"ecomm/internal/domain"
	"ecomm/internal/dto"
	"ecomm/internal/repository"
	notification "ecomm/pkg/Notification"
	"errors"
	"log"
	// "strconv"
	"fmt"
	"time"

	"ecomm/internal/helper"
)

type UserService struct {
	// Add necessary fields like repository, logger, etc.
	Repo   repository.UserRepository
	Catalog   repository.CatalogRepository
	Auth   helper.Auth
	Config configs.AppConfig
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

// check if user is verified or not
func (us UserService) IsUserVerified(id uint) bool {
	// Implement the logic to check if the user is verified.
	currentUser, err := us.Repo.FindUserByID(id)
	return err == nil && currentUser.Verified
}

func (us UserService) GetVerificationCode(e domain.User) error {
	//if user already verified
	if us.IsUserVerified(e.ID) {
		return errors.New("user already verified")
	}
	//generate verification code
	code, err := us.Auth.GenerateCode()
	if err != nil {
		return err
	}
	//update the user with latest verification code
	user := domain.User{
		Expiry: time.Now().Add(30 * time.Minute),
		Code:   code,
	}
	_, err = us.Repo.UpdateUser(e.ID, user)
	if err != nil {
		return errors.New("failed to update user with verification code")
	}
	//grab the user donot care about error bcoz we have already checked the user exist or not
	user, _ = us.Repo.FindUserByID(e.ID)

	//send sms or email to user with the code
	notificationClient := notification.NewNotificationClient(us.Config)
	msg := fmt.Sprintf("Your verification code is: %v", code)
	err = notificationClient.SendSMS(user.Phone, msg)
	if err != nil {
		return errors.New("Error on sending sms")
	}
	//return verification code
	return nil
}

func (us UserService) VerifyCode(id uint, code int) error {
	if us.IsUserVerified(id) {
		log.Println("verified...")
		return errors.New("user already verified")
	}
	user, err := us.Repo.FindUserByID(id)
	if err != nil {
		return err
	}
	if user.Code != code {
		return errors.New("verification code doesn't match")
	}
	if time.Now().After(user.Expiry) {
		return errors.New("Verification code expired")
	}
	//if everything is fine update the user as verified
	updateUser := domain.User{
		Verified: true,
	}
	_, err = us.Repo.UpdateUser(id, updateUser)

	if err != nil {
		return errors.New("Unable to Verify User")
	}

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

func (us UserService) BecomeSeller(id uint, input dto.SellerInput) (string, error) {
	//find existing user
	user, _ := us.Repo.FindUserByID(id)
	//return already a seller
	if user.UserType == domain.SELLER {
		return "", errors.New("user already a seller")
	}

	//update the user type to seller
	seller, err := us.Repo.UpdateUser(id, domain.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Phone:     input.PhoneNumber,
		UserType:  domain.SELLER,
	})

	if err != nil {
		return "", errors.New("failed to update user to seller")
	}

	//generate and return the token
	token, err := us.Auth.GenerateToken(seller.ID, seller.Email, seller.UserType)

	//create bank account information for the seller
	account := domain.BankAccount{
		UserID:      id,
		BankAccount: input.BankAccountNumber,
		IFSCCode:    input.IFSCCode,
		PaymentType: input.PaymentType,
	}
	err = us.Repo.CreateBankAccount(account)

	return token, err
}

// return bunch of card item so return slice of interface
func (us UserService) FindCart(id uint) ([]domain.Cart, error) {
	cartItems, err := us.Repo.FindCartItems(id)
	log.Printf("error %v",err)

	return cartItems, err
}

// input as product properties and user info to update the cart
func (us UserService) CreateCart(input dto.CreateCartRequest, user domain.User) ([]domain.Cart, error) {
	//check if cart exist for the user
	cart, _ := us.Repo.FindCartItem(user.ID, input.ProductID)

	if cart.ID > 0 {
		//=>delete the cart items
		if input.Quantity < 1 {
			err := us.Repo.DeleteCartById(cart.ID)
			log.Printf("errors on deleting cart items %v", err)
			if err != nil {
				return nil, errors.New("errors on deleting cart items")
			}
		} else {
			//=> update the cart items
			cart.Qty = int(input.Quantity)
			err := us.Repo.UpdateCart(cart)
			log.Printf("errors on updating cart items %v", err)
			if err != nil {
				return nil, errors.New("errors on updating cart items")
			}
		}
		// return us.Repo.FindCartItems(user.ID)
	} else {
		//=> create new cart items
		product, _:= us.Catalog.FindProductByID(int(input.ProductID))
		if product.ID>0 {
			return nil, errors.New("product not found to create cart item")
		}
		err:= us.Repo.CreateCart(domain.Cart{
			UserId:    user.ID,
			ProductId: input.ProductID,
			Name:      product.Name,
			ImageUrl: product.ImageUrl,
			Qty:	   int(input.Quantity),
			Price:    product.Price,
			SellerId: uint(product.UserId),
		})
		if err != nil {
			return nil, errors.New("errors on creating cart items")
		}

	return us.Repo.FindCartItems(user.ID)
}

// just find the user whether the user have cart or not
func (us UserService) CreateOrder(user *domain.User) (int, error) {

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
