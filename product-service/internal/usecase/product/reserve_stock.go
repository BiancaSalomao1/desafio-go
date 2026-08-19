package product

import (
	"product-service/internal/repository"
)

type StockItem struct {
	ProductID string
	Quantity  int
}

type ReserveStockUseCase struct {
	productRepository repository.ProductRepository
}

func NewReserveStockUseCase(
	productRepository repository.ProductRepository,
) *ReserveStockUseCase {
	return &ReserveStockUseCase{
		productRepository: productRepository,
	}
}

func (uc *ReserveStockUseCase) Execute(
	items []StockItem,
) error {
	for _, item := range items {
		product, err := uc.productRepository.FindByID(item.ProductID)
		if err != nil {
			return err
		}

		if err := product.ReduceStock(item.Quantity); err != nil {
			return err
		}

		if err := uc.productRepository.Update(product); err != nil {
			return err
		}
	}

	return nil
}
