package service

import (
	"context"
	"github.com/birukbelay/item/entity"
	"github.com/birukbelay/item/packages/items"
	"github.com/birukbelay/item/utils/helpers"
)

// CategoriesService implements items.CategoriesService interface
type CategoriesService struct {
	categoriesRepo items.CategoriesRepository
}

// NewCategoriesService will create new CategoriesService object
func NewCategoriesService(CatRepo items.CategoriesRepository) items.CategoriesService {
	return &CategoriesService{categoriesRepo: CatRepo}
}


// Categories returns list of Categoriess
func (cs *CategoriesService) Categories(ctx context.Context, limit int , offset string) ([]entity.Categories, []error) {

	Categoriess, errs := cs.categoriesRepo.Categories(ctx, limit, offset)

	if len(errs) > 0 {
		return nil, errs
	}

	return Categoriess, nil
}

// StoreCategories persists new categories information
func (cs *CategoriesService) StoreCategories(ctx context.Context, categories *entity.Categories) (*entity.Categories, []error) {

	cat, errs := cs.categoriesRepo.StoreCategories(ctx, categories)

	if len(errs) > 0 {
		return nil, errs
	}

	return cat, nil
}

// Category returns a categories object with a given id
func (cs *CategoriesService) Category(ctx context.Context, id string) (*entity.Categories, []error) {

	c, errs := cs.categoriesRepo.Category(ctx, id)

	if len(errs) > 0 {
		return c, errs
	}

	return c, nil
}

// UpdateCategories updates a cateogory with new data
func (cs *CategoriesService) UpdateCategories(ctx context.Context, categories *entity.Categories) (*entity.Categories, []error) {

	helpers.LogTrace("id", categories.ID)
	cat, errs := cs.categoriesRepo.UpdateCategories(ctx, categories)

	if len(errs) > 0 {
		return nil, errs
	}

	return cat, nil
}

// DeleteCategories delete a categories by its id
func (cs *CategoriesService) DeleteCategories(ctx context.Context, id string) (*entity.Categories, []error) {

	cat, errs := cs.categoriesRepo.DeleteCategories(ctx, id)

	if len(errs) > 0 {
		return nil, errs
	}

	return cat, nil
}


