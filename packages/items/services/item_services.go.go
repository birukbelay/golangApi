package service

import (
	 "github.com/birukbelay/item/entity"
	"github.com/birukbelay/item/packages/items"
)

// ItemService implements items.ItemService interface
type ItemServices struct {
	itemRepo items.ItemRepository
}

// NewItemService returns new ItemService object
func NewItemService(itemRepository items.ItemRepository) items.ItemService {
	return &ItemServices{itemRepo: itemRepository}
}

func (is *ItemServices) ItemsByFilter(limit int, offsetValue, offsetType, categories, brand, stype, sort string, sort_way int, min_price, max_price int) ([]entity.Item, string, string, []error) {
	itms,_,_, errs := is.itemRepo.ItemsByFilter(limit , offsetValue, offsetType, categories, brand, stype, sort , sort_way , min_price, max_price )
	if len(errs) > 0 {
		return nil,"","", errs
	}
	return itms,"","", errs
}

// Items returns all stored item items items
func (is *ItemServices) Items(limit, offset int) ([]entity.Item, []error) {
	itms, errs := is.itemRepo.Items(limit, offset)
	if len(errs) > 0 {
		return nil, errs
	}
	return itms, errs
}



// Item retrieves a item items items by its id
func (is *ItemServices) Item(id string) (*entity.Item, []error) {
	itm, errs := is.itemRepo.Item(id)
	if len(errs) > 0 {
		return nil, errs
	}
	return itm, errs
}

// UpdateItem updates a given item items items
func (is *ItemServices) UpdateItem(item *entity.Item) (*entity.Item, []error) {
	itm, errs := is.itemRepo.UpdateItem(item)
	if len(errs) > 0 {
		return nil, errs
	}
	return itm, errs
}

// DeleteItem deletes a given item items items
func (is *ItemServices) DeleteItem(id string) (int64, []error) {
	itm, errs := is.itemRepo.DeleteItem(id)
	if len(errs) > 0 {
		return 0, errs
	}
	return itm, errs
}

// StoreItem stores a given item items items
func (is *ItemServices) StoreItem(item *entity.Item) (*entity.Item, []error) {
	itm, errs := is.itemRepo.StoreItem(item)
	if len(errs) > 0 {
		return nil, errs
	}
	return itm, errs
}
