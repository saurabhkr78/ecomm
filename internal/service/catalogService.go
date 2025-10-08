package service

import (
	"ecomm/configs"
	"ecomm/internal/domain"
	"ecomm/internal/dto"
	"ecomm/internal/helper"
	"ecomm/internal/repository"
	"errors"
	"strconv"
)

type CatalogService struct {
	// Add necessary fields like repository, logger, etc.
	Repo   repository.CatalogRepository
	Auth   helper.Auth
	Config configs.AppConfig
}

func (s CatalogService) CreateCategory(input dto.CreateCategoryRequest) error {
	err := s.Repo.CreateCategory(domain.Category{
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

	return &UpdatedCategory, err
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
func (s CatalogService) GetCategories() (*[]domain.Category, error) {
	categories, err := s.Repo.FindCategories()
	if err != nil {
		return nil, errors.New("Categories doesnot exist")
	}
	return &categories, err
}

func (s CatalogService) GetCategory(id uint) (*domain.Category, error) {
	category, err := s.Repo.FindCategoryByID(id)
	if err != nil {
		return nil, errors.New("Category doesnot exist")
	}
	return &category, err
}
