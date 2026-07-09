package product

/*
struct CreateProductRequest

Responsabilidades:

- receber os dados para criação de um produto.

Campos:
- name
- price
- stock
*/

type CreateProductRequest struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Stock int     `json:"stock"`
}
