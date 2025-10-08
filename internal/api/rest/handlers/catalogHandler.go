package handlers

import (
	"ecomm/internal/api/rest"
	"ecomm/internal/dto"
	"strconv"

	"ecomm/internal/repository"
	"ecomm/internal/service"

	"github.com/gofiber/fiber/v3"
)

type CatalogHandler struct {
	//catelog service
	svc service.CatalogService
}

// here we need to accept something in our setup user routes function which is `app *fiber.App
//
//	so that we can register our routes with the fiber app instance
//
// so we need to have another struct which will have the fiber app instance
// and then we can call the setup user routes function using that struct instance so we created httpHandler.go
func SetupCatalogRoutes(rh *rest.RestHandler) {
	//now here we can grab the fiber app spinning in server.js using rh.App

	app := rh.App

	//so,in future when we gonna create the instance of user service and inject to handler
	svc := service.CatalogService{
		Repo:   repository.NewCatalogRepository(rh.DB),
		Auth:   rh.Auth,
		Config: rh.Config,
	}

	handler := &CatalogHandler{
		svc: svc,
	}

	//grouping the routes

	/*
		Here public end points will be listing products and categories
		and private end points will be managing products and categories
	*/
	// ---------- Public endpoints ----------
	app.Get("/products", handler.GetProducts)           // List all products
	app.Get("/products/:id", handler.GetProduct)        // Get a single product
	app.Get("/categories", handler.GetCategories)       // List all categories
	app.Get("/categories/:id", handler.GetCategoryById) // Get a single category)
	// ---------- Private endpoints ----------
	//
	selRoutes := app.Group("/seller", rh.Auth.AuthorizeSeller) //seller specific routes

	//catagories
	selRoutes.Post("/categories", handler.CreateCategory)       // Add a new category
	selRoutes.Patch("/categories/:id", handler.EditCategory)    // Update a category
	selRoutes.Delete("/categories/:id", handler.DeleteCategory) // Delete a category

	//products
	selRoutes.Post("/products", handler.CreateProducts)
	selRoutes.Get("/products", handler.GetProducts)
	selRoutes.Get("/products/:id", handler.GetProduct)
	selRoutes.Put("/products/:id", handler.EditProducts)
	selRoutes.Patch("/products/:id/stocks", handler.UpdateStocks)
	selRoutes.Delete("/products/:id", handler.DeleteProducts)
}

//public route fxns

func (ch CatalogHandler) GetCategories(ctx fiber.Ctx) error {

	categories, err := ch.svc.GetCategories()
	if err != nil {
		return rest.ErrorMessage(ctx, 404, err)
	}

	return rest.SuccessResponse(ctx, "Categories ", categories)
}

func (ch CatalogHandler) GetCategoryById(ctx fiber.Ctx) error {

	id, _ := strconv.Atoi(ctx.Params("id"))
	category, err := ch.svc.GetCategory(uint(id))
	if err != nil {
		return rest.ErrorMessage(ctx, 404, err)
	}
	return rest.SuccessResponse(ctx, "Category ", category)
}

func (ch CatalogHandler) GetProduct(ctx fiber.Ctx) error {
	return rest.SuccessResponse(ctx, "Get Product by id endpoint", nil)
}

func (ch CatalogHandler) GetProducts(ctx fiber.Ctx) error {
	return rest.SuccessResponse(ctx, "Get Products endpoint", nil)
}

//private routes fxn for seller

func (ch CatalogHandler) CreateCategory(ctx fiber.Ctx) error {

	req := dto.CreateCategoryRequest{}
	err := ctx.Bind().Body(&req)
	if err != nil {
		return rest.BadRequestError(ctx, "Create category request is not valid")
	}

	err = ch.svc.CreateCategory(req)
	if err != nil {
		return rest.InternalError(ctx, err)
	}

	return rest.SuccessResponse(ctx, "Category Created Successfully", nil)
}

//edit category

func (ch CatalogHandler) EditCategory(ctx fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))
	req := dto.CreateCategoryRequest{}
	err := ctx.Bind().Body(&req)

	if err != nil {
		return rest.BadRequestError(ctx, "edit category request is not valid")
	}

	updatedCategory, err := ch.svc.EditCategory(uint(id), req)
	if err != nil {
		return rest.InternalError(ctx, err)
	}
	return rest.SuccessResponse(ctx, "Edit Category", updatedCategory)
}

// func (ch CatalogHandler) EditCategory(ctx fiber.Ctx) error {
// 	// 1. Get the category ID from URL param
// 	id, err := strconv.Atoi(ctx.Params("id"))
// 	if err != nil {
// 		return rest.BadRequestError(ctx, "invalid category ID")
// 	}

// 	// 2. Create a struct to hold the JSON body
// 	var req dto.CreateCategoryRequest

// 	// 3. Bind the JSON body into the struct
// 	if err := ctx.Bind().Body(&req); err != nil {
// 		return rest.BadRequestError(ctx, "edit category request is not valid")
// 	}

// 	// 4. Call the service to update the category
// 	updatedCategory, err := ch.svc.EditCategory(uint(id), req)
// 	if err != nil {
// 		return rest.InternalError(ctx, err)
// 	}

// 	// 5. Return success response
// 	return rest.SuccessResponse(ctx, "Edit Category", updatedCategory)
// }

func (ch *CatalogHandler) DeleteCategory(ctx fiber.Ctx) error {
	id, _ := strconv.Atoi(ctx.Params("id"))
	err := ch.svc.DeleteCategory(uint(id))
	if err != nil {
		return rest.InternalError(ctx, err)
	}
	return rest.SuccessResponse(ctx, " Category Deleted Successfully", nil)
}

// product
func (ch CatalogHandler) CreateProducts(ctx fiber.Ctx) error {
	return rest.SuccessResponse(ctx, "Create Products endpoint", nil)
}
func (ch CatalogHandler) EditProducts(ctx fiber.Ctx) error {
	return rest.SuccessResponse(ctx, "Edit Products endpoint", nil)

}

func (ch CatalogHandler) UpdateStocks(ctx fiber.Ctx) error {
	return rest.SuccessResponse(ctx, "Update Stocks endpoint", nil)
}

func (ch CatalogHandler) DeleteProducts(ctx fiber.Ctx) error {
	return rest.SuccessResponse(ctx, "Delete Products endpoint", nil)
}
