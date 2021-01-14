package item

import (
	"context"
	"main/utils/tracing"

	"go.opentelemetry.io/otel/label"
)

// Service Item Storage Service
type Service interface {
	GetAllItems(ctx context.Context) ([]Item, error)
	GetItemsForBarcode(ctx context.Context, barcode string) ([]Item, error)
	Update(ctx context.Context, item Item) (string, error)
}

type service struct {
	dao Dao
}

// NewService Create Item Service
func NewService(dao Dao) Service {
	return &service{
		dao: dao,
	}
}

func (s *service) GetAllItems(ctx context.Context) ([]Item, error) {
	sc, span := tracing.CreateSpan(ctx, "item-service", "get-all")
	defer span.End()

	return s.dao.GetAllItems(sc), nil
}

func (s *service) GetItemsForBarcode(ctx context.Context, barcode string) ([]Item, error) {
	sc, span := tracing.CreateSpan(ctx, "item-service", "get-barcode")
	defer span.End()

	return s.dao.GetItemsForBarcode(sc, barcode)
}

func (s *service) Update(ctx context.Context, item Item) (string, error) {
	sc, span := tracing.CreateSpan(ctx, "item-service", "update")
	defer span.End()

	span.SetAttributes(
		label.String("barcode", item.Barcode),
		label.String("name", item.Name),
		label.Int("quantity", item.Quantity),
	)

	err := s.dao.Update(sc, item)
	if err != nil {
		return "", err
	}
	return "Successfully Updated", nil
}
