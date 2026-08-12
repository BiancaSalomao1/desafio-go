package product

/*
struct UpdateProductRequest

Responsabilidades:
- receber os dados para atualização de um produto.

Campos:
- name
- price
- stock

Métodos:
- nenhum
*/

type UpdateProductRequest struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Stock int     `json:"stock"`
}
