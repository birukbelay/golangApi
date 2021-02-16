package items

import "github.com/birukbelay/item/entity"


// ItemRepository specifies item menu items related database operations
type ItemRepository interface {
	Items(limit, offset int) ([]entity.Item, []error)
	Item(id string) (*entity.Item, []error)
	UpdateItem(menu *entity.Item) (*entity.Item, []error)
	DeleteItem(id string) (int64, []error)
	StoreItem(item *entity.Item) (*entity.Item, []error)

	ItemsByFilter(limit int, offsetValue, offsetType,
		categories, brand, stype ,
		sort string, sort_way int,
		min_price, max_price int) ([]entity.Item, string, string, []error)

}

// CategoriesRepository specifies item menu categories database operations
type CategoriesRepository interface {
	Categories(limit int , offset string) ([]entity.Categories, []error)
	Category(id string) (*entity.Categories, []error)
	UpdateCategories(categories *entity.Categories) (*entity.Categories, []error)
	DeleteCategories(id string) (*entity.Categories, []error)
	StoreCategories(categories *entity.Categories) (*entity.Categories, []error)
	//ItemsInCategories(categories *entity.Category, limit, offset int) ([]entity.Item, []error)
}
