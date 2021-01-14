package item

import (
	"context"
	"errors"
	"main/utils/tracing"

	"go.opentelemetry.io/otel/label"
)

type memStore struct {
	barcodeToNameToQuantity map[string]map[string]int
}

// NewDaoInMemory Create In Memory Storage
func NewDaoInMemory(ctx context.Context, initalizeData bool) Dao {
	sc, span := tracing.CreateSpan(ctx, "item-dao", "create")
	defer span.End()

	span.SetAttributes(label.Bool("initializeData", initalizeData))

	m := &memStore{
		barcodeToNameToQuantity: make(map[string]map[string]int),
	}
	span.AddEvent("Created Store")

	if initalizeData {
		m.initialize(sc)
		span.AddEvent("Data Injected")
	}

	return m
}

func (m *memStore) initialize(ctx context.Context) {
	sc, span := tracing.CreateSpan(ctx, "item-dao", "init")
	defer span.End()

	err := m.Update(sc, Item{Barcode: "1", Name: "Toilet-Paper", Quantity: 1})
	if err != nil {
		span.RecordError(err)
		return
	}
	err = m.Update(sc, Item{Barcode: "2", Name: "Toilet-Paper", Quantity: 1})
	if err != nil {
		span.RecordError(err)
		return
	}
	err = m.Update(sc, Item{Barcode: "1", Name: "Paper-Towels", Quantity: 1})
	if err != nil {
		span.RecordError(err)
		return
	}
	err = m.Update(sc, Item{Barcode: "3", Name: "Ziploc-Gallon", Quantity: 1})
	if err != nil {
		span.RecordError(err)
		return
	}
	err = m.Update(sc, Item{Barcode: "4", Name: "Ziploc-Quart", Quantity: 3})
	if err != nil {
		span.RecordError(err)
		return
	}
	span.AddEvent("Mock Data Initalized")
}

func (m *memStore) GetAllItems(ctx context.Context) (items []Item) {
	_, span := tracing.CreateSpan(ctx, "item-dao", "get-all")
	defer span.End()

	for barcode, nTq := range m.barcodeToNameToQuantity {
		for name, quantity := range nTq {
			items = append(items, Item{Barcode: barcode, Name: name, Quantity: quantity})
		}
	}
	SortItems(items)

	return
}

func (m *memStore) GetItemsForBarcode(ctx context.Context, barcode string) (items []Item, err error) {
	_, span := tracing.CreateSpan(ctx, "item-dao", "get-barcode")
	defer span.End()

	nTq, hasBarcode := m.barcodeToNameToQuantity[barcode]
	if !hasBarcode {
		err = errors.New("barcode not found")
		span.RecordError(err)
		return
	}

	for name, quantity := range nTq {
		items = append(items, Item{Barcode: barcode, Name: name, Quantity: quantity})
	}
	SortItems(items)

	return
}

func (m *memStore) Update(ctx context.Context, item Item) error {
	_, span := tracing.CreateSpan(ctx, "item-dao", "adjust-quantity")
	defer span.End()

	if _, hasBarcode := m.barcodeToNameToQuantity[item.Barcode]; !hasBarcode {
		m.barcodeToNameToQuantity[item.Barcode] = make(map[string]int)
	}

	m.barcodeToNameToQuantity[item.Barcode][item.Name] = item.Quantity

	return nil
}
