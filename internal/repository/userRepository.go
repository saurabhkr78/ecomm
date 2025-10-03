// how this repository will interact with the database
/*
while initiating the UserService.go somehow we need to inject the userRepository there so that userService can interact with the database through userRepository
so while UserService functions are called it has to call the DB inorder to perform db operations

*/
/*
Why people use interfaces here

Testing: you can easily make a “fake” repository or service for unit tests.

Swapping implementations: you can later switch from DB to API or a different DB without changing the rest of your code.

Decoupling: your controller only knows about the interface, not about the actual struct.
*/

package repository

import (
	"ecomm/internal/domain"
	"errors"
	"log"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepository interface {
	//CRUD operations
	CreateUser(usr domain.User) (domain.User, error)
	FindUser(email string) (domain.User, error)
	FindUserByID(id uint) (domain.User, error)
	UpdateUser(id uint, usr domain.User) (domain.User, error)

	//more function will come here like delete user, get all users etc
}
type userRepository struct {
	//db connection will come here
	db *gorm.DB
}

// since we cannot create object of interface so we need to create object of struct which implements the interface
// constructor function to return the object of userRepository struct
/*
while initalizing user routes in UserHandler we cannot directly pass the interface so we need to create object of struct as struct declared over there is a private repository so,  implements the interface and pass it to UserHandler

*/
/* now userRepositry implemented all the function if i delete a fxn then it gonna give me error beacuse userReposittry is not the kind of userrepository interrface anymore*/
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

// func (r userRepository) CreateUser(usr domain.User) (domain.User, error) {
// 	err := r.db.Create(&usr).Error
// 	if err != nil {
// 		log.Printf("error while creating user: %v\n", err)
// 		return domain.User{}, errors.New("could not create user")
// 	}
// 	return usr, err
// }

func (r userRepository) CreateUser(usr domain.User) (domain.User, error) {
	// The Create method will populate the usr object with the new ID.
	err := r.db.Create(&usr).Error
	if err != nil {
		log.Printf("error while creating user: %v\n", err)
		return domain.User{}, errors.New("could not create user")
	}

	// After creating, we fetch the record again using the new ID.
	// This is the most reliable way to ensure all fields, including
	// database-level defaults like 'user_type', are correctly loaded.
	var createdUser domain.User
	if err := r.db.First(&createdUser, usr.ID).Error; err != nil {
		log.Printf("failed to fetch user after creation: %v", err)
		return domain.User{}, errors.New("could not retrieve user after creation")
	}

	// Return the fully populated user object.
	return createdUser, nil
}

func (r userRepository) FindUser(email string) (domain.User, error) {

	var user domain.User
	err := r.db.First(&user, "email = ?", email).Error
	if err != nil {
		log.Printf("error while finding user by email: %v\n", err)
		return domain.User{}, errors.New("user not exists")
	}
	return user, nil
}

func (r userRepository) FindUserByID(id uint) (domain.User, error) {
	var user domain.User
	err := r.db.First(&user, id).Error
	if err != nil {
		log.Printf("error while finding user by id: %v\n", err)
		return domain.User{}, errors.New("user not exists")
	}
	return user, nil
}

func (r userRepository) UpdateUser(id uint, usr domain.User) (domain.User, error) {
	var user domain.User
	err := r.db.Model(&user).Clauses(clause.Returning{}).Where("id = ?", id).Updates(usr).Error
	if err != nil {
		log.Printf("error while updating user: %v\n", err)
		return domain.User{}, errors.New("failed to update user")
	}
	return user, nil
}
