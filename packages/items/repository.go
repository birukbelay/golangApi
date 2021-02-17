package items

import (
	"context"
	"github.com/birukbelay/item/entity"
)


// ItemRepository specifies item menu items related database operations
type ItemRepository interface {
	Items(ctx context.Context,limit, offset int) ([]entity.Item, []error)
	Item(ctx context.Context, id string) (*entity.Item, []error)
	UpdateItem(ctx context.Context, menu *entity.Item) (*entity.Item, []error)
	DeleteItem(ctx context.Context, id string) (int64, []error)
	StoreItem(ctx context.Context, item *entity.Item) (*entity.Item, []error)

	ItemsByFilter(ctx context.Context, limit int, offsetValue, offsetType,
		categories, brand, stype ,
		sort string, sort_way int,
		min_price, max_price int) ([]entity.Item, string, string, []error)

}

// CategoriesRepository specifies item menu categories database operations
type CategoriesRepository interface {
	Categories(ctx context.Context, limit int , offset string) ([]entity.Categories, []error)
	Category(ctx context.Context, id string) (*entity.Categories, []error)
	UpdateCategories(ctx context.Context, categories *entity.Categories) (*entity.Categories, []error)
	DeleteCategories(ctx context.Context, id string) (*entity.Categories, []error)
	StoreCategories(ctx context.Context, categories *entity.Categories) (*entity.Categories, []error)
	//ItemsInCategories(categories *entity.Category, limit, offset int) ([]entity.Item, []error)
}
