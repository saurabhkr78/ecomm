package payment

import (
	"fmt"

	"github.com/stripe/stripe-go/v83"
	"github.com/stripe/stripe-go/v83/checkout/session"
	"log"
)

type PaymentClient interface {
	CreatePayment(amount float64, userId uint, orderId string) (*stripe.CheckoutSession, error)
	GetPaymentStatus(pId string) (*stripe.CheckoutSession, error)
}

type Payment struct {
	stripeSecretKey string
	successUrl      string
	cancelUrl       string
}

// func (p Payment) CreatePayment(amount float64, userId uint, orderId string) (*stripe.CheckoutSession, error) {
// 	stripe.Key = p.stripeSecretKey
// 	amountInCents := float64(amount * 100)

// 	params := &stripe.CheckoutSessionParams{
// 		PaymentMethodTypes: stripe.StringSlice([]string{"card", "UPI", "Netbanking", "wallets", "emi", "paylater", "paytm"}),
// 		LineItems: []*stripe.CheckoutSessionLineItemParams{
// 			{
// 				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
// 					UnitAmount: stripe.Int64(int64(amountInCents)),
// 					Currency:   stripe.String("usd"),
// 					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
// 						Name: stripe.String("Electronics Order Payment"),
// 					},
// 				},
// 				Quantity: stripe.Int64(1),
// 			},
// 		},
// 		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
// 		SuccessURL: stripe.String(p.successUrl),
// 		CancelURL:  stripe.String(p.cancelUrl),
// 	}
// 	params.AddMetadata("user_id", fmt.Sprintf("%d", userId))
// 	params.AddMetadata("order_id", fmt.Sprintf("%s", orderId))

//		sess, err := session.New(params)
//		if err != nil {
//			log.Printf("session creation error: %v", err)
//			return nil, fmt.Errorf("failed to create stripe checkout session: %w", err)
//		}
//		return sess, nil
//	}
func (p Payment) CreatePayment(amount float64, userId uint, orderId string) (*stripe.CheckoutSession, error) {
	stripe.Key = p.stripeSecretKey
	amountInCents := float64(amount * 100)

	params := &stripe.CheckoutSessionParams{
		// Use only valid Stripe payment method types
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					UnitAmount: stripe.Int64(int64(amountInCents)),
					Currency:   stripe.String("usd"),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String("Electronics Order Payment"),
					},
				},
				Quantity: stripe.Int64(1),
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModePayment)),
		SuccessURL: stripe.String(p.successUrl),
		CancelURL:  stripe.String(p.cancelUrl),
	}
	params.AddMetadata("user_id", fmt.Sprintf("%d", userId))
	params.AddMetadata("order_id", fmt.Sprintf("%s", orderId))

	sess, err := session.New(params)
	if err != nil {
		log.Printf("session creation error: %v", err)
		return nil, fmt.Errorf("failed to create stripe checkout session: %w", err)
	}
	return sess, nil
}

func (p Payment) GetPaymentStatus(pId string) (*stripe.CheckoutSession, error) {
	stripe.Key = p.stripeSecretKey

	sess, err := session.Get(pId, nil)
	if err != nil {
		log.Printf("failed to get payment status: %v", err)
		return nil, fmt.Errorf("failed to get payment status: %w", err)
	}
	return sess, nil
}

// NewPaymentClient creates a new instance of PaymentClient
// with the provided Stripe secret key, success URL, and cancel URL.
// It returns a PaymentClient interface.

//this is a constructor function for PaymentClient

func NewPaymentClient(stripeSecretKey string) PaymentClient {

	return &Payment{
		stripeSecretKey: stripeSecretKey,
	}
}
