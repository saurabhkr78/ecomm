package repository

import (
	"ecomm/internal/domain"
	"errors"
	"log"

	"gorm.io/gorm"
)

type CatalogRepository interface {
	CreateCategory(e domain.Category) error
	FindCategories() ([]domain.Category, error)
	FindCategoryByID(id uint) (domain.Category, error)
	EditCategory(e domain.Category) (domain.Category, error)
	DeleteCategory(id uint) error
}
type catalogRepository struct {
	//db connection will come here
	db *gorm.DB
}

func (c catalogRepository) CreateCategory(e domain.Category) error {
	err := c.db.Create(&e).Error

	if err != nil {
		log.Printf("db_error %v", err)
		return errors.New("Failed to create category")
	}
	return nil
}

func (c catalogRepository) FindCategories() ([]domain.Category, error) {
	var categories []domain.Category
	err := c.db.Find(&categories).Error

	if err != nil {
		return nil, err
	}
	return categories, nil
}

func (c catalogRepository) FindCategoryByID(id uint) (domain.Category, error) {
	var category domain.Category
	err := c.db.First(&category, id).Error
	if err != nil {
		log.Printf("db_error %v", err)
		return domain.Category{}, errors.New("category doesn't exist")
	}
	return category, nil
}

func (c catalogRepository) EditCategory(e domain.Category) (domain.Category, error) {
	err := c.db.Save(&e).Error
	if err != nil {
		log.Printf("db_error %v", err)
		return domain.Category{}, errors.New("Failed to update category")
	}
	return e, nil
}
func (c catalogRepository) DeleteCategory(id uint) error {
	err := c.db.Delete(&domain.Category{}, id).Error
	if err != nil {
		log.Printf("db_error %v", err)
		return errors.New("Failed to delete category")
	}
	return nil
}

// since we cannot create object of interface so we need to create object of struct which implements the interface
// constructor function to return the object of userRepository struct
/*
while initalizing user routes in UserHandler we cannot directly pass the interface so we need to create object of struct as struct declared over there is a private repository so,  implements the interface and pass it to UserHandler

*/
/* now catalogRepository implemented all the function if i delete a fxn then it gonna give me error beacuse catalogRepository is not the kind of catalogRepository interrface anymore*/
func NewCatalogRepository(db *gorm.DB) CatalogRepository {
	return &catalogRepository{
		db: db,
	}
}
