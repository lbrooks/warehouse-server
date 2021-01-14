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
	AdjustQuantity(ctx context.Context, barcode, name string, count int) (string, error)
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

func (s *service) AdjustQuantity(ctx context.Context, barcode, name string, count int) (string, error) {
	sc, span := tracing.CreateSpan(ctx, "item-service", "adjust-quantity")
	defer span.End()

	span.SetAttributes(
		label.String("barcode", barcode),
		label.String("name", name),
		label.Int("count", count),
	)

	err := s.dao.AdjustQuantity(sc, barcode, name, count)
	if err != nil {
		return "", err
	}
	return "Successfully Updated", nil
}
