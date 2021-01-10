package item

import (
	"context"
	"errors"
	"main/utils/slice"
	"main/utils/tracing"
	"sort"

	"go.opentelemetry.io/otel/label"
)

type memStore struct {
	barcodeToItemName  map[string][]string
	itemNameToQuantity map[string]int
}

// NewDaoInMemory Create In Memory Storage
func NewDaoInMemory(ctx context.Context, initalizeData bool) Dao {
	sc, span := tracing.CreateSpan(ctx, "item-dao", "create")
	defer span.End()

	span.SetAttributes(label.Bool("initializeData", initalizeData))

	m := &memStore{
		barcodeToItemName:  make(map[string][]string),
		itemNameToQuantity: make(map[string]int),
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

	m.AdjustQuantity(sc, "1", "Toilet-Paper", 1)
	m.AdjustQuantity(sc, "2", "Toilet-Paper", 1)
	m.AdjustQuantity(sc, "1", "Paper-Towels", 1)
	m.AdjustQuantity(sc, "3", "Ziploc-Gallon", 1)
	m.AdjustQuantity(sc, "4", "Ziploc-Quart", 1)
	span.AddEvent("Mock Data Initalized")
}

func (m *memStore) GetAllItems(ctx context.Context) (items []Item) {
	_, span := tracing.CreateSpan(ctx, "item-dao", "get-all")
	defer span.End()

	for k, v := range m.itemNameToQuantity {
		items = append(items, Item{Name: k, Quantity: v})
	}
	span.AddEvent("Added Inventory")

	sort.Slice(items, func(i, j int) bool {
		return items[i].Name < items[j].Name
	})

	return
}

func (m *memStore) GetItemsForBarcode(ctx context.Context, barcode string) (items []Item, err error) {
	_, span := tracing.CreateSpan(ctx, "item-dao", "get-barcode")
	defer span.End()

	names, hasBarcode := m.barcodeToItemName[barcode]
	if !hasBarcode {
		err = errors.New("barcode not found")
		span.RecordError(err)
		return
	}

	for _, name := range names {
		quantity, hasName := m.itemNameToQuantity[name]
		if !hasName {
			err = errors.New("name not found")
			span.RecordError(err)
			return
		}

		items = append(items, Item{Name: name, Quantity: quantity})
	}

	return
}

func (m *memStore) AdjustQuantity(ctx context.Context, barcode, name string, add int) error {
	_, span := tracing.CreateSpan(ctx, "item-dao", "adjust-quantity")
	defer span.End()

	if name == "" {
		err := errors.New("missing a name")
		span.RecordError(err)
		return err
	}

	if barcode != "" {
		if names, hasBarcode := m.barcodeToItemName[barcode]; !hasBarcode {
			m.barcodeToItemName[barcode] = []string{name}
		} else if !slice.Contains(names, name) {
			m.barcodeToItemName[barcode] = append(m.barcodeToItemName[barcode], name)
		}
	}

	if _, hasQuantity := m.itemNameToQuantity[name]; !hasQuantity {
		m.itemNameToQuantity[name] = 0
	}

	m.itemNameToQuantity[name] += add

	if m.itemNameToQuantity[name] < 0 {
		m.itemNameToQuantity[name] = 0
	}

	return nil
}
