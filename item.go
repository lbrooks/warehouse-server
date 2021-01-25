package server

import (
	"context"
	"sort"
)

// Item in my inventory
type Item struct {
	Barcode  string `form:"barcode" json:"barcode"`
	Brand    string `form:"brand" json:"brand"`
	Name     string `form:"name" json:"name"`
	Quantity int    `form:"quantity" json:"quantity" binding:"min=0"`
}

// SortItems sort a slice of items
func SortItems(items []*Item) {
	sort.Slice(items, func(i, j int) bool {
		isEqual := items[i].Name == items[j].Name
		if !isEqual {
			return items[i].Name < items[j].Name
		}
		isEqual = items[i].Brand == items[j].Brand
		if !isEqual {
			return items[i].Brand < items[j].Brand
		}
		isEqual = items[i].Barcode == items[j].Barcode
		if !isEqual {
			return items[i].Barcode < items[j].Barcode
		}
		return items[i].Quantity < items[j].Quantity
	})
}

// ItemService Item Service
type ItemService interface {
	GetCounts(ctx context.Context) (map[string]int, error)
	Search(ctx context.Context, item Item) ([]*Item, error)
	Update(ctx context.Context, item Item) (string, error)
}
