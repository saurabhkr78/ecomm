package service

import (
	"ecomm/configs"
	"ecomm/internal/domain"
	"ecomm/internal/dto"
	"ecomm/internal/helper"
	"ecomm/internal/repository"
	notification "ecomm/pkg/Notification"
	"errors"
	"log"
	// "strconv"
	"fmt"
	"time"
)

type UserService struct {
	// Add necessary fields like repository, logger, etc.
	Repo    repository.UserRepository
	Catalog repository.CatalogRepository
	Auth    helper.Auth
	Config  configs.AppConfig
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

func (us UserService) CreateProfile(id uint, input dto.ProfileInput) error {
	//find the user first
	user, err := us.Repo.FindUserByID(id)
	if err != nil {
		return err
	}
	if input.FirstName != "" {
		user.FirstName = input.FirstName
	}
	if input.LastName != "" {
		user.LastName = input.LastName
	}

	_, err = us.Repo.UpdateUser(id, user)

	//create address
	address := domain.Address{
		AddressLine1: input.AddressInput.AddressLine1,
		AddressLine2: input.AddressInput.AddressLine2,
		Street:       input.AddressInput.Street,
		City:         input.AddressInput.City,
		State:        input.AddressInput.State,
		Country:      input.AddressInput.Country,
		PinCode:      input.AddressInput.PinCode,
		UserID:       id,
	}
	err = us.Repo.CreateProfile(address)
	if err != nil {
		return err
	}

	return nil
}

// find the user by id and return the user profile but sometime
// we are using pointer bcoz at any point of time we need to edit the specific profile so taht why we are returning kind of pointer
func (us UserService) GetProfile(id uint) (*domain.User, error) {
	user, err := us.Repo.FindUserByID(id)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
func (us UserService) UpdateProfile(id uint, input dto.ProfileInput) error {
	//find the user first
	user, err := us.Repo.FindUserByID(id)
	if err != nil {
		return err
	}
	if input.FirstName != "" {
		user.FirstName = input.FirstName
	}
	if input.LastName != "" {
		user.LastName = input.LastName
	}

	_, err = us.Repo.UpdateUser(id, user)

	address := domain.Address{
		AddressLine1: input.AddressInput.AddressLine1,
		AddressLine2: input.AddressInput.AddressLine2,
		Street:       input.AddressInput.Street,
		City:         input.AddressInput.City,
		State:        input.AddressInput.State,
		Country:      input.AddressInput.Country,
		PinCode:      input.AddressInput.PinCode,
		UserID:       id,
	}

	err = us.Repo.UpdateProfile(address)
	if err != nil {
		return err
	}

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
	log.Printf("error %v", err)

	return cartItems, err
}

// func (us UserService) CreateCart(input dto.CreateCartRequest, u domain.User) ([]domain.Cart, error) {
// 	// check if the cart is Exist
// 	cart, _ := us.Repo.FindCartItem(u.ID, input.ProductID)

// 	if cart.ID > 0 {
// 		if input.ProductID == 0 {
// 			return nil, errors.New("please provide a valid product id")
// 		}
// 		//  => delete the cart item
// 		if input.Quantity < 1 {
// 			err := us.Repo.DeleteCartById(cart.ID)
// 			if err != nil {
// 				log.Printf("Error on deleting cart item %v", err)
// 				return nil, errors.New("error on deleting cart item")
// 			}
// 		} else {
// 			//  => update the cart item
// 			cart.Qty = int(input.Quantity)
// 			err := us.Repo.UpdateCart(cart)
// 			if err != nil {
// 				// log error
// 				return nil, errors.New("error on updating cart item")
// 			}
// 		}

// 	} else {
// 		// check if product exist
// 		product, _ := us.Catalog.FindProductByID(input.ProductID)
// 		if product.ID < 1 {
// 			return nil, errors.New("product not found to create cart item")
// 		}
// 		// create cart

// 		err := us.Repo.CreateCart(domain.Cart{
// 			UserId:    u.ID,
// 			ProductId: input.ProductID,
// 			Name:      product.Name,
// 			ImageUrl:  product.ImageUrl,
// 			Qty:       int(input.Quantity),
// 			Price:     product.Price,
// 			SellerId:  uint(product.UserId),
// 		})

// 		if err != nil {
// 			return nil, errors.New("error on creating cart item")
// 		}
// 	}

// 	return us.Repo.FindCartItems(u.ID)

// }

func (us UserService) CreateCart(input dto.CreateCartRequest, u domain.User) ([]domain.Cart, error) {
	// check if the cart is Exist
	cart, _ := us.Repo.FindCartItem(u.ID, input.ProductID)

	// check if product exist
	product, err := us.Catalog.FindProductByID(input.ProductID)
	if err != nil {
		return nil, errors.New("product not found")
	}

	if cart.ID > 0 {
		if input.ProductID == 0 {
			return nil, errors.New("please provide a valid product id")
		}
		//  => delete the cart item
		if input.Quantity < 1 {
			err := us.Repo.DeleteCartById(cart.ID)
			if err != nil {
				log.Printf("Error on deleting cart item %v", err)
				return nil, errors.New("error on deleting cart item")
			}
		} else {
			//  => update the cart item
			cart.Qty = cart.Qty + int(input.Quantity)
			err := us.Repo.UpdateCart(cart)
			if err != nil {
				// log error
				return nil, errors.New("error on updating cart item")
			}
		}

	} else {
		if product.ID < 1 {
			return nil, errors.New("product not found to create cart item")
		}
		// create cart

		err := us.Repo.CreateCart(domain.Cart{
			UserId:    u.ID,
			ProductId: input.ProductID,
			Name:      product.Name,
			ImageUrl:  product.ImageUrl,
			Qty:       int(input.Quantity),
			Price:     product.Price,
			SellerId:  uint(product.UserId),
		})

		if err != nil {
			return nil, errors.New("error on creating cart item")
		}
	}

	return us.Repo.FindCartItems(u.ID)

}

// just find the user whether the user have cart or not
func (us UserService) CreateOrder(user domain.User) (int, error) {
	//find cart items of that user
	cartItems, err := us.Repo.FindCartItems(user.ID)
	if err != nil {
		return 0, errors.New("error on finding cart items for the user")
	}
	if len(cartItems) == 0 {
		return 0, errors.New("Cart is empty cannot create the order")
	}

	//find successful payment of that user
	paymentId := "PAY12345" //dummy payment id
	txnId := "TXN12345"     //dummy txn id
	orderRef := helper.RandomNumbers(8)

	//create Order with Generated OrderNumber
	var totalAmount float64
	var orderItems []domain.OrderItem
	for _, item := range cartItems {
		totalAmount += float64(item.Qty) * item.Price
		orderItem := domain.OrderItem{
			ProductId: item.ProductId,
			Name:      item.Name,
			Quantity:  item.Qty,
			Price:     item.Price,
			ImageUrl:  item.ImageUrl,
			SellerId:  item.SellerId,
		}
		orderItems = append(orderItems, orderItem)
	}
	order := domain.Order{
		UserID:         user.ID,
		PaymentId:      paymentId,
		TransactionId:  txnId,
		OrderRefNumber: uint(orderRef),
		TotalAmount:    totalAmount,
		Items:          orderItems,
	}
	err := us.Repo.CreateOrder(order)

	//send email to user with order details

	//remove cart items from the cart

	//return order number

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
