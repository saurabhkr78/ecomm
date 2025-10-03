package helper

import (
	"ecomm/internal/domain"
	"errors"
	"strings"
	"time"

	"fmt"
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	Secret string
}

/*
 all the fxns are receiver fxns of Auth struct
so we can create an instance of Auth struct and call these fxns using that instance
this is similar to how we created UserHandler struct and called its receiver fxns using its instance
this is called method set in go
*/

func SetUpAuth(s string) Auth {
	return Auth{
		Secret: s,
	}
}

func (a Auth) CreateHashedPassword(password string) (string, error) {
	if len(password) < 6 {
		return "", errors.New("Password is too short")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		//this are internal errors so shouldbe logged using some logger
		return "", errors.New("Failed to hash password")
	}
	return string(hashedPassword), nil
}

func (a Auth) GenerateToken(id uint, email string, role string) (string, error) {

	if id == 0 || email == "" || role == "" {
		return "", errors.New("Invalid data to generate token")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": id,
		"email":   email,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24 * 30).Unix(), // Token expires in 72 hours
	})

	tokenString, err := token.SignedString([]byte(a.Secret))
	if err != nil {
		return "", errors.New("Failed to generate token")
	}
	return tokenString, nil
}

func (a Auth) VerifyPassword(plainpassword string, hashedpassword string) error {

	if len(plainpassword) < 6 {
		return errors.New("Password is too short")
	}
	//Grab the error if any
	err := bcrypt.CompareHashAndPassword([]byte(hashedpassword), []byte(plainpassword))
	if err != nil {
		return errors.New("Password does not match")
	}
	return nil
}

func (a Auth) VerifyToken(token string) (domain.User, error) {
	// Split token into two parts: "Bearer <token>"
	tokenArr := strings.Split(token, " ")

	if len(tokenArr) != 2 {
		return domain.User{}, nil
	}

	tokenStr := tokenArr[1]
	// Check the first part is "Bearer", not the second
	if tokenArr[0] != "Bearer" {
		return domain.User{}, errors.New("invalid token format")
	}

	// Parse the JWT token
	parsedToken, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		// Check signing method
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(a.Secret), nil
	})
	if err != nil {
		return domain.User{}, fmt.Errorf("invalid token: %w", err)
	}

	// Validate token and extract claims
	if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
		// Check expiration
		if exp, ok := claims["exp"].(float64); ok {
			if time.Now().Unix() > int64(exp) {
				return domain.User{}, errors.New("token expired")
			}
		} else {
			return domain.User{}, errors.New("token missing expiration")
		}

		// Extract user information
		user := domain.User{}
		if userID, ok := claims["user_id"].(float64); ok {
			user.ID = uint(userID)
		} else {
			return domain.User{}, errors.New("token missing user_id")
		}

		if email, ok := claims["email"].(string); ok {
			user.Email = email
		} else {
			return domain.User{}, errors.New("token missing email")
		}

		if role, ok := claims["role"].(string); ok {
			user.UserType = role
		} else {
			return domain.User{}, errors.New("token missing role")
		}

		return user, nil
	}

	return domain.User{}, errors.New("token verification failed")
}

func (a Auth) Authorize(ctx fiber.Ctx) error {
	headers := ctx.GetReqHeaders()
	authHeaders, exists := headers["Authorization"]

	if !exists || len(authHeaders) == 0 {
		return ctx.Status(401).JSON(fiber.Map{
			"message": "unauthorized",
			"reason":  "missing authorization header",
		})
	}

	authHeader := authHeaders[0] // Take the first Authorization header
	user, err := a.VerifyToken(authHeader)

	if err != nil {
		return ctx.Status(401).JSON(fiber.Map{
			"message": "unauthorized",
			"reason":  err.Error(),
		})
	}

	if user.ID == 0 {
		return ctx.Status(401).JSON(fiber.Map{
			"message": "unauthorized",
			"reason":  "invalid user",
		})
	}
	//locals automatically gets deleted after request is processed
	//locals are used to pass data between middlewares and handlers
	//you can store any data you want in locals
	ctx.Locals("user", user)
	return ctx.Next()
}

func (a Auth) GetCurrentUser(ctx fiber.Ctx) domain.User {

	user := ctx.Locals("user")

	return user.(domain.User)
}
