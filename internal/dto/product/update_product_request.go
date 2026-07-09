package product

/*
struct UpdateProductRequest

Responsabilidades:

- receber os dados para atualização de um produto.

Campos:
- id
- name
- price
- stock
*/

type UpdateProductRequest struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Stock int     `json:"stock"`
}
