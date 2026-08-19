package product

import (
	"product-service/internal/repository"
)

type ReleaseStockUseCase struct {
	productRepository repository.ProductRepository
}

func NewReleaseStockUseCase(
	productRepository repository.ProductRepository,
) *ReleaseStockUseCase {
	return &ReleaseStockUseCase{
		productRepository: productRepository,
	}
}

func (uc *ReleaseStockUseCase) Execute(
	items []StockItem,
) error {
	for _, item := range items {
		product, err := uc.productRepository.FindByID(item.ProductID)
		if err != nil {
			return err
		}

		if err := product.IncreaseStock(item.Quantity); err != nil {
			return err
		}

		if err := uc.productRepository.Update(product); err != nil {
			return err
		}
	}

	return nil
}
