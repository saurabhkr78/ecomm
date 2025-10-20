package service

import (
	"ecomm/configs"
	"ecomm/internal/domain"
	"ecomm/internal/dto"
	"ecomm/internal/helper"
	"ecomm/internal/repository"
	"errors"
	"fmt"
	"strconv"
)

type CatalogService struct {
	// Add necessary fields like repository, logger, etc.
	Repo   repository.CatalogRepository
	Auth   helper.Auth
	Config configs.AppConfig
}

func (s CatalogService) CreateCategory(input dto.CreateCategoryRequest) error {
	err := s.Repo.CreateCategory(&domain.Category{
		Name:         input.Name,
		ImageUrl:     input.ImageUrl,
		DisplayOrder: input.DisplayOrder,
	})
	return err
}

func (s CatalogService) EditCategory(id uint, input dto.CreateCategoryRequest) (*domain.Category, error) {
	existingCategory, err := s.Repo.FindCategoryByID(id)
	if err != nil {
		return nil, errors.New("Category does not exist")
	}
	if len(input.Name) > 0 {
		existingCategory.Name = input.Name
	}
	if input.ParentId != nil {
		// Convert the uint ParentId to a string to match the domain model
		existingCategory.ParentId = strconv.FormatUint(uint64(*input.ParentId), 10)
	}
	if len(input.ImageUrl) > 0 {
		existingCategory.ImageUrl = input.ImageUrl
	}
	if input.DisplayOrder >= 0 {
		existingCategory.DisplayOrder = input.DisplayOrder
	}
	UpdatedCategory, err := s.Repo.EditCategory(existingCategory)

	return UpdatedCategory, err
}

func (s CatalogService) DeleteCategory(id uint) error {
	err := s.Repo.DeleteCategory(id)
	if err != nil {

		//log the error
		return errors.New("Category does not exist to delete")
	}
	return nil
}

// errror check on both fxn is needed
func (s CatalogService) GetCategories() ([]*domain.Category, error) {
	categories, err := s.Repo.FindCategories()
	if err != nil {
		return nil, errors.New("Categories doesnot exist")
	}
	return categories, err
}

func (s CatalogService) GetCategory(id uint) (*domain.Category, error) {
	category, err := s.Repo.FindCategoryByID(id)
	if err != nil {
		return nil, errors.New("Category doesnot exist")
	}
	return category, err
}

// Product service functions
func (s CatalogService) CreateProduct(input dto.CreateProductRequest, user domain.User) error {
	err := s.Repo.CreateProduct(&domain.Product{
		Name:        input.Name,
		Description: input.Description,
		Price:       float64(input.Price),
		CategoryId:  fmt.Sprintf("%d", input.CategoryId),
		ImageUrl:    input.ImageUrl,
		Stock:       int(input.Stock),
		UserId:      int(user.ID),
	})
	return err
}
func (s CatalogService) EditProduct(id uint, input dto.CreateProductRequest, user domain.User) (*domain.Product, error) {
	existingProduct, err := s.Repo.FindProductByID(id)
	if err != nil {
		return nil, errors.New("Product does not exist")
	}

	//verify product owner
	if existingProduct.UserId != int(user.ID) {
		return nil, errors.New("Unauthorized: You do not own this product or right to edit the product")
	}

	if len(input.Name) > 0 {
		existingProduct.Name = input.Name
	}
	if len(input.Description) > 0 {
		existingProduct.Description = input.Description
	}
	// if input.ParentId > 0 {
	// 	existingProduct.ParentId = input.ParentId
	// }
	if len(input.ImageUrl) > 0 {
		existingProduct.ImageUrl = input.ImageUrl
	}
	// if input.DisplayOrder >= 0 {
	// 	existingProduct.DisplayOrder = input.DisplayOrder
	// }
	if input.Price > 0 {
		existingProduct.Price = float64(input.Price)
	}
	if input.CategoryId > 0 {
		existingProduct.CategoryId = fmt.Sprintf("%d", input.CategoryId)
	}

	UpdatedProduct, err := s.Repo.EditProduct(existingProduct)

	return UpdatedProduct, err
}
func (s CatalogService) DeleteProduct(id uint, user domain.User) error {
	existingProduct, err := s.Repo.FindProductByID(id)
	if err != nil {
		return errors.New("Product does not exist to delete")
	}
	//verify product owner
	if existingProduct.UserId != int(user.ID) {
		return errors.New("Unauthorized: You do not own this product or have the manage right to delete the product")
	}
	err = s.Repo.DeleteProduct(existingProduct.ID)
	if err != nil {
		return errors.New("Failed to delete product")
	}
	return nil
}
func (s CatalogService) GetProducts() ([]*domain.Product, error) {
	products, err := s.Repo.FindProducts()
	if err != nil {
		return nil, errors.New("Products doesnot exist")
	}
	return products, err
}
func (s CatalogService) GetProductById(id uint) (*domain.Product, error) {
	product, err := s.Repo.FindProductByID(id)
	if err != nil {
		return nil, errors.New("Product doesnot exist")
	}
	return product, err
}

func (s CatalogService) GetSellerProducts(id uint) ([]*domain.Product, error) {
	products, err := s.Repo.FindSellerProducts(id)
	if err != nil {
		return nil, errors.New("Products doesnot exist for this seller")
	}
	return products, err
}
func (s CatalogService) UpdateProductStock(id uint, stock int, user domain.User) (*domain.Product, error) {
	existingProduct, err := s.Repo.FindProductByID(id)
	if err != nil {
		return nil, errors.New("Product does not exist")
	}
	//verify product owner
	if existingProduct.UserId != int(user.ID) {
		return nil, errors.New("Unauthorized: You do not own this product or right to edit the product stock")
	}
	existingProduct.Stock = stock
	updatedProduct, err := s.Repo.EditProduct(existingProduct)
	if err != nil {
		return nil, errors.New("Failed to update product stock")
	}
	return updatedProduct, nil

}
