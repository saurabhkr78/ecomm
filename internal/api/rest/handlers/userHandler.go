/*
/ UserHandler is responsible for handling user-related HTTP requests.
// It depends on the UserService to perform business logic operations.
type UserHandler struct {
    svc service.UserService // Service layer dependency injected into the handler
}

// SetupUserRoutes registers all user-related routes with the Fiber app.
// It takes a RestHandler, which holds the Fiber app instance, so routes can be attached.
func SetupUserRoutes(rh *rest.RestHandler) {
    // Grab the Fiber app instance from RestHandler (this is your running server).
    app := rh.App

    // Create an instance of the UserService.
    // This service contains the business logic for user operations.
    svc := service.UserService{}

    // Create a new UserHandler and inject the UserService into it.
    // This is an example of dependency injection, which makes testing and maintenance easier.
    handler := &UserHandler{svc: svc}

    // Now you would register your routes with the Fiber app.
    // Example:
    // app.Get("/users", handler.GetUsers)
    // app.Post("/users", handler.CreateUser)




*/

package handlers

import (
	"ecomm/internal/api/rest"
	"ecomm/internal/dto"
	"ecomm/internal/service"
	"github.com/gofiber/fiber/v3"
	"net/http"
)

// since all the receiver function will have handler function so we can create a struct type
// once the handler instance is created we can call the receiver function using that instance
// when calling any endpoint our user handler will be able to respond accordingly as part of the API calls
type UserHandler struct {
	svc service.UserService
}

// here we need to accept something in our setup user routes function which is `app *fiber.App
//
//	so that we can register our routes with the fiber app instance
//
// so we need to have another struct which will have the fiber app instance
// and then we can call the setup user routes function using that struct instance so we created httpHandler.go
func SetupUserRoutes(rh *rest.RestHandler) {
	//now here we can grab the fiber app spinning in server.js using rh.App

	app := rh.App

	//so,in future when we gonna create the instance of user service and inject to handler
	svc := service.UserService{}
	handler := &UserHandler{svc: svc}

	// ---------- Public endpoints ----------
	app.Post("/register", handler.Register) // User signup
	app.Post("/login", handler.Login)       // User login

	// ---------- Private endpoints ----------
	app.Get("/verify", handler.GetVerificationCode) // Fetch verification code (e.g., for OTP)
	app.Post("/verify", handler.Verify)             // Submit verification

	app.Post("/profile", handler.CreateProfile) // Create/update profile
	app.Get("/profile", handler.GetProfile)     // Fetch profile

	app.Post("/cart", handler.AddToCart) // Add/update cart
	app.Get("/cart", handler.GetCart)    // View cart

	app.Get("/order", handler.GetOrders)    // Fetch all orders
	app.Get("/order/:id", handler.GetOrder) // Fetch specific order

	app.Post("/become-seller", handler.BecomeSeller) // Upgrade to seller account
}

// receiver function that wil accept the user handler instance
// receiver function is only appalicable when we are going to create one kind of instance of specific handler

// ---------------- HANDLER METHODS ----------------

// Register handles new user registration
func (uh *UserHandler) Register(ctx fiber.Ctx) error {
	//some dto validation
	//call the service layer
	//return

	//to create user
	user := dto.UserSignup{}
	err := ctx.Bind().Body(&user)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "please provide valid inputs",
		})
	}
	token, err := uh.svc.Signup(user)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "error on signup",
		})
	}
	return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": token})
}

// Login handles user login
func (uh *UserHandler) Login(ctx fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": "login"})
}

// GetVerificationCode returns a verification code
func (uh *UserHandler) GetVerificationCode(ctx fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": "get verification code"})
}

// Verify submits a verification code
func (uh *UserHandler) Verify(ctx fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": "verify"})
}

// CreateProfile creates or updates a user profile
func (uh *UserHandler) CreateProfile(ctx fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": "profile created/updated"})
}

// GetProfile fetches the user profile
func (uh *UserHandler) GetProfile(ctx fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": "profile"})
}

// AddToCart adds items to the user's cart
func (uh *UserHandler) AddToCart(ctx fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": "item added to cart"})
}

// GetCart fetches the user's cart items
func (uh *UserHandler) GetCart(ctx fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": "cart items"})
}

// GetOrders fetches all user orders
func (uh *UserHandler) GetOrders(ctx fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": "all orders"})
}

// GetOrder fetches a specific order by ID
func (uh *UserHandler) GetOrder(ctx fiber.Ctx) error {
	orderID := ctx.Params("id", "all")
	return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": "order", "order_id": orderID})
}

// BecomeSeller upgrades a user to a seller account
func (uh *UserHandler) BecomeSeller(ctx fiber.Ctx) error {
	return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": "become seller"})
}
