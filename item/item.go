package item

import "sort"

// Item in my inventory
type Item struct {
	Barcode  string `form:"barcode" binding:"required"`
	Name     string `form:"name" binding:"required"`
	Quantity int    `form:"quantity" binding:"required,min=0"`
}

// SortItems sort a slice of items
func SortItems(items []Item) {
	sort.Slice(items, func(i, j int) bool {
		isEqual := items[i].Name == items[j].Name
		if !isEqual {
			return items[i].Name < items[j].Name
		}
		isEqual = items[i].Quantity == items[j].Quantity
		if !isEqual {
			return items[i].Quantity < items[j].Quantity
		}
		return items[i].Barcode < items[j].Barcode
	})
}
