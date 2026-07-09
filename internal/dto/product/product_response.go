package product

/*
struct ProductResponse

Responsabilidades:

- retornar os dados de um produto.

Campos:
- id
- name
- price
- stock
*/

type ProductResponse struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Stock int     `json:"stock"`
}
