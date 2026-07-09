package domain

/*
struct Product

- identificar o produto;
- armazenar nome;
- armazenar preço;
- controlar estoque.

Métodos:
- construtor NewProduct()
- ReduceStock()
- IncreaseStock()
- Validate()
*/

type Product struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Stock int     `json:"stock"`
}

// NewProduct cria uma nova instância de Product.
func NewProduct(
	id string,
	name string,
	price float64,
	stock int,
) *Product {
	return &Product{
		ID:    id,
		Name:  name,
		Price: price,
		Stock: stock,
	}
}

func (p *Product) Validate() error {

	if p.Name == "" {
		return ErrProductNameRequired
	}

	if p.Price <= 0 {
		return ErrInvalidPrice
	}

	if p.Stock < 0 {
		return ErrInvalidStock
	}

	return nil
}

func (p *Product) ReduceStock(quantity int) error {

	if quantity <= 0 {
		return ErrInvalidQuantity
	}

	if quantity > p.Stock {
		return ErrInsufficientStock
	}

	p.Stock -= quantity

	return nil
}

// IncreaseStock devolve produtos ao estoque.
func (p *Product) IncreaseStock(quantity int) error {

	if quantity <= 0 {
		return ErrInvalidQuantity
	}

	p.Stock += quantity

	return nil
}
