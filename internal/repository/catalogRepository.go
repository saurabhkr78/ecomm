package repository

import (
	"ecomm/internal/domain"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
)

type CatalogRepository interface {
	CreateCategory(e domain.Category) error
	FindCategories() ([]domain.Category, error)
	FindCategoryByID(id uint) (domain.Category, error)
	EditCategory(e domain.Category) (domain.Category, error)
	DeleteCategory(id uint) error

	// Add other necessary methods for catalog management
	CreateProduct(e *domain.Product) error
	FindProducts() ([]*domain.Product, error)
	FindProductByID(id uint) (*domain.Product, error)
	FindSellerProducts(sellerID uint) ([]*domain.Product, error)
	EditProduct(e *domain.Product) (*domain.Product, error)
	DeleteProduct(id uint) error
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

//product repository functions

func (c *catalogRepository) CreateProduct(e *domain.Product) error {
	err := c.db.Model(&domain.Product{}).Create(e).Error
	if err != nil {
		log.Printf("db_error %v", err)
		return errors.New("failed to create product")
	}
	return nil
}

func (c *catalogRepository) FindProducts() ([]*domain.Product, error) {
	var products []*domain.Product
	err := c.db.Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (c *catalogRepository) FindProductByID(id uint) (*domain.Product, error) {
	var product *domain.Product
	err := c.db.First(&product, id).Error
	if err != nil {
		log.Printf("db_error %v", err)
		return nil, errors.New("product doesn't exist")
	}
	return product, nil
}

func (c *catalogRepository) FindSellerProducts(id uint) ([]*domain.Product, error) {
	var products []*domain.Product
	err := c.db.Where("user_id = ?", id).Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (c *catalogRepository) EditProduct(e *domain.Product) (*domain.Product, error) {
	err := c.db.Save(&e).Error
	if err != nil {
		log.Printf("db_error %v", err)
		return nil, errors.New("failed to update product")
	}
	return e, nil
}

func (c *catalogRepository) DeleteProduct(id uint) error {
	// err := c.db.Delete(&domain.Product{}, id).Error
	// if err != nil {
	// 	return fmt.Errorf("failed to delete product with id %d: %w", id, err)
	// }

	// return nil
	result := c.db.Delete(&domain.Product{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("no product found with id or already deleted %d", id)
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
