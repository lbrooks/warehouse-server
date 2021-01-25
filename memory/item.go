package memory

import (
	"context"
	"fmt"

	"main/inventory"
	"main/utils/booleans"
	"main/utils/tracing"

	"go.opentelemetry.io/otel/label"
)

type memStore struct {
	items []*inventory.Item
}

func (m *memStore) findItemsMatching(ctx context.Context, barcode, name, brand string) (items []*inventory.Item) {
	_, span := tracing.CreateSpan(ctx, "item.daomem.memStore", "findItemsMatching")
	defer span.End()

	ignoreEmpty := booleans.IsBlank(barcode) && booleans.IsBlank(name) && booleans.IsBlank(brand)
	items = make([]*inventory.Item, 0)
	for _, v := range m.items {
		barcodeMatch := booleans.IsBlank(barcode) || v.Barcode == barcode
		nameMatch := booleans.IsBlank(name) || v.Name == name
		brandMatch := booleans.IsBlank(brand) || v.Brand == brand

		if barcodeMatch && nameMatch && brandMatch {
			if !ignoreEmpty || v.Quantity > 0 {
				items = append(items, v)
			}
		}
	}
	return
}

func (m *memStore) initialize(ctx context.Context) {
	sc, span := tracing.CreateSpan(ctx, "item-dao", "initialize")
	defer span.End()

	err := m.Update(sc, inventory.Item{Barcode: "1", Name: "Toilet Paper", Brand: "Charmin", Quantity: 1})
	if err != nil {
		span.RecordError(err)
		return
	}
	err = m.Update(sc, inventory.Item{Barcode: "2", Name: "Toilet Paper", Brand: "Sandpaper", Quantity: 1})
	if err != nil {
		span.RecordError(err)
		return
	}
	err = m.Update(sc, inventory.Item{Barcode: "1", Name: "Paper Towels", Brand: "Bounty", Quantity: 1})
	if err != nil {
		span.RecordError(err)
		return
	}
	err = m.Update(sc, inventory.Item{Barcode: "3", Name: "Gallon Bag", Brand: "Ziploc", Quantity: 1})
	if err != nil {
		span.RecordError(err)
		return
	}
	err = m.Update(sc, inventory.Item{Barcode: "4", Name: "Quart Bag", Brand: "Ziploc", Quantity: 3})
	if err != nil {
		span.RecordError(err)
		return
	}
	span.AddEvent("Mock Data Initalized")
}

// NewDaoInMemory Create In Memory Storage
func NewDaoInMemory(ctx context.Context, initalizeData bool) Dao {
	sc, span := tracing.CreateSpan(ctx, "item.dao.daomem", "NewDaoInMemory")
	defer span.End()

	span.SetAttributes(label.Bool("initializeData", initalizeData))

	m := &memStore{
		items: make([]*inventory.Item, 0),
	}
	span.AddEvent("Created Store")

	if initalizeData {
		m.initialize(sc)
		span.AddEvent("Data Injected")
	}

	return m
}

func (m *memStore) GetCounts(ctx context.Context) map[string]int {
	_, span := tracing.CreateSpan(ctx, "item-dao", "GetCounts")
	defer span.End()

	counts := make(map[string]int)
	for _, i := range m.items {
		if _, found := counts[i.Name]; !found {
			counts[i.Name] = i.Quantity
		} else {
			counts[i.Name] += i.Quantity
		}
	}
	return counts
}

func (m *memStore) Search(ctx context.Context, item inventory.Item) (items []*inventory.Item) {
	_, span := tracing.CreateSpan(ctx, "item-dao", "Search")
	defer span.End()

	items = m.findItemsMatching(ctx, item.Barcode, item.Name, item.Brand)
	inventory.SortItems(items)

	return
}

func (m *memStore) Update(ctx context.Context, item inventory.Item) error {
	_, span := tracing.CreateSpan(ctx, "item-dao", "Update")
	defer span.End()

	span.SetAttributes(
		label.String("barcode", item.Barcode),
		label.String("name", item.Name),
		label.Int("quantity", item.Quantity),
	)

	matching := m.findItemsMatching(ctx, item.Barcode, item.Name, item.Brand)
	if len(matching) == 0 {
		m.items = append(m.items, &item)
	} else if len(matching) > 1 {
		err := fmt.Errorf("Multiple Items Matched")
		span.RecordError(err)
		return err
	} else {
		matching[0].Quantity = item.Quantity
		if matching[0].Quantity < 0 {
			matching[0].Quantity = 0
		}
	}

	return nil
}
