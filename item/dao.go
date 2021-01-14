package item

import "context"

// Dao Item Storage Data Access Layer
type Dao interface {
	GetAllItems(ctx context.Context) []Item
	GetItemsForBarcode(ctx context.Context, barcode string) ([]Item, error)
	Update(ctx context.Context, item Item) error
}
