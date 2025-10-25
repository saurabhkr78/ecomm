package handlers

import (
	"ecomm/configs"
	"ecomm/internal/api/rest"
	"ecomm/internal/helper"
	"ecomm/internal/repository"
	"ecomm/internal/service"
	"ecomm/pkg/payment"
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
	"net/http"
)

type TransactionHandler struct {
	Svc           service.TransactionService
	UserSvc       service.UserService
	PaymentClient payment.PaymentClient
	Config        configs.AppConfig
}

func initializeTransactionService(db *gorm.DB, auth helper.Auth) service.TransactionService {
	return service.TransactionService{
		Repo: repository.NewTransactionRepository(db),
		Auth: auth,
	}
}

func SetupTransactionRoutes(as *rest.RestHandler) {

	app := as.App
	svc := initializeTransactionService(as.DB, as.Auth)
	useSvc := service.UserService{
		Repo:    repository.NewUserRepository(as.DB),
		Catalog: repository.NewCatalogRepository(as.DB),
		Auth:    as.Auth,
		Config:  as.Config,
	}

	handler := TransactionHandler{
		Svc:           svc,
		PaymentClient: as.Pc,
		UserSvc:       useSvc,
		Config:        as.Config,
	}

	secRoute := app.Group("/buyer", as.Auth.Authorize)
	secRoute.Get("/payment", handler.MakePayment)
	secRoute.Get("/verify", handler.VerifyPayment)

	sellerRoute := app.Group("/seller", as.Auth.AuthorizeSeller)
	sellerRoute.Get("/orders", handler.GetOrders)
	sellerRoute.Get("/orders/:id", handler.GetOrderDetails)
}

func (h *TransactionHandler) MakePayment(ctx fiber.Ctx) error {

	//1. grab authorized user
	user := h.Svc.Auth.GetCurrentUser(ctx)

	//2.check if payment session active then return the payment url
	activePayment, err := h.Svc.GetActivePayment(user.ID)
	if err != nil {
		// Handle error, maybe log and return internal server error
		return rest.InternalError(ctx, err)
	}
	if activePayment != nil && activePayment.ID > 0 {
		return ctx.Status(http.StatusOK).JSON(&fiber.Map{
			"message":     "existing payment",
			"payment_url": activePayment.PaymentUrl,
		})
	}
	//3.Calll user service get cart data to aggregate the total amount and collect payment
	_, amount, err := h.UserSvc.FindCart(user.ID)

	orderId, err := helper.RandomNumbers(8)
	if err != nil {
		return rest.InternalError(ctx, errors.New("error generating order id"))
	}

	// 4. create a new payment session

	sessionResult, err := h.PaymentClient.CreatePayment(amount, user.ID, orderId)

	// 5. store payment session data to create to store payment
	err = h.Svc.StoreCreatedPayment(user.ID, sessionResult, amount, orderId)
	if err != nil {
		return ctx.Status(400).JSON(err)
	}
	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message":     "create payment",
		"result":      sessionResult,
		"payment_url": sessionResult.URL,
	})
}

func (h *TransactionHandler) VerifyPayment(ctx fiber.Ctx) error {

	// grab authorized user
	user := h.Svc.Auth.GetCurrentUser(ctx)

	// do we have active payment session to verify?
	activePayment, err := h.Svc.GetActivePayment(user.ID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return rest.InternalError(ctx, err)
	}
	if activePayment == nil || activePayment.ID == 0 {
		return ctx.Status(400).JSON(errors.New("no active payment exist"))
	}

	// fetch payment status from stripe
	paymentRes, err := h.PaymentClient.GetPaymentStatus(activePayment.PaymentId)
	paymentJson, _ := json.Marshal(paymentRes)
	paymentLogs := string(paymentJson)
	paymentStatus := "failed"

	// if payment then create order
	if paymentRes.Status == "succeeded" {
		// create Order
		paymentStatus = "success"
		err = h.UserSvc.CreateOrder(user.ID, activePayment.OrderId, activePayment.PaymentId, activePayment.Amount)
	}

	if err != nil {
		return rest.InternalError(ctx, err)
	}

	// update payment status
	h.Svc.UpdatePayment(user.ID, paymentStatus, paymentLogs)

	return ctx.Status(http.StatusOK).JSON(&fiber.Map{
		"message":  "create payment",
		"response": paymentRes,
	})
}

func (h *TransactionHandler) GetOrders(ctx fiber.Ctx) error {
	return ctx.Status(200).JSON("success")
}

func (h *TransactionHandler) GetOrderDetails(ctx fiber.Ctx) error {
	return ctx.Status(200).JSON("success")
}
