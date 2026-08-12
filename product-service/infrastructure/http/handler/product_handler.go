package handler

/*
struct ProductHandler

Responsabilidades:
- receber requisições HTTP relacionadas a produtos;
- converter DTOs para entidades de domínio;
- executar os casos de uso;
- retornar respostas HTTP em JSON.
- atualizar um produto.

Campos:
- createProductUseCase
- getProductUseCase
- listProductsUseCase
- updateProductUseCase
-deleteProductUseCase

Métodos:
- NewProductHandler()
- Create()
- GetByID()
- List()
- Update()
- Delete()

*/

import (
	"net/http"

	"product-service/infrastructure/http/httpx"

	productdto "product-service/internal/dto/product"
	"product-service/internal/mapper"
)

type ProductHandler struct {
	createProductUseCase CreateProductUseCase
	getProductUseCase    GetProductUseCase
	listProductsUseCase  ListProductsUseCase
	updateProductUseCase UpdateProductUseCase
	deleteProductUseCase DeleteProductUseCase
}

func NewProductHandler(
	createProductUseCase CreateProductUseCase,
	getProductUseCase GetProductUseCase,
	listProductsUseCase ListProductsUseCase,
	updateProductUseCase UpdateProductUseCase,
	deleteProductUseCase DeleteProductUseCase,
) *ProductHandler {
	return &ProductHandler{
		createProductUseCase: createProductUseCase,
		getProductUseCase:    getProductUseCase,
		listProductsUseCase:  listProductsUseCase,
		updateProductUseCase: updateProductUseCase,
		deleteProductUseCase: deleteProductUseCase,
	}
}

// Create
//
// @Summary Criar Produto
// @Tags Products
// @Accept json
// @Produce json
// @Param request body product.CreateProductRequest true "Produto"
// @Success 201 {object} product.ProductResponse
// @Security BearerAuth
// @Router /products [post]
func (h *ProductHandler) Create(
	w http.ResponseWriter,
	r *http.Request,
) {

	var request productdto.CreateProductRequest

	if err := httpx.ReadJSON(
		r,
		&request,
	); err != nil {

		httpx.WriteError(
			w,
			http.StatusBadRequest,
			err,
		)

		return
	}

	product := mapper.ToProduct(request)

	if err := h.createProductUseCase.Execute(product); err != nil {

		httpx.WriteError(
			w,
			http.StatusBadRequest,
			err,
		)

		return
	}

	httpx.WriteJSON(
		w,
		http.StatusCreated,
		mapper.ToProductResponse(product),
	)
}

// List
//
// @Summary Listar Produtos
// @Tags Products
// @Produce json
// @Success 200 {array} product.ProductResponse
// @Security BearerAuth
// @Router /products [get]
func (h *ProductHandler) List(
	w http.ResponseWriter,
	r *http.Request,
) {

	products, err := h.listProductsUseCase.Execute()

	if err != nil {

		httpx.WriteError(
			w,
			http.StatusInternalServerError,
			err,
		)

		return
	}

	httpx.WriteJSON(
		w,
		http.StatusOK,
		mapper.ToProductResponseList(products),
	)
}

// GetByID
//
// @Summary Buscar Produto
// @Tags Products
// @Produce json
// @Param id path string true "ID"
// @Success 200 {object} product.ProductResponse
// @Security BearerAuth
// @Router /products/{id} [get]
func (h *ProductHandler) GetByID(
	w http.ResponseWriter,
	r *http.Request,
) {

	id := r.PathValue("id")

	product, err := h.getProductUseCase.Execute(id)

	if err != nil {

		httpx.WriteError(
			w,
			http.StatusNotFound,
			err,
		)

		return
	}

	httpx.WriteJSON(
		w,
		http.StatusOK,
		mapper.ToProductResponse(product),
	)
}

// Update
//
// @Summary Atualizar Produto
// @Tags Products
// @Accept json
// @Produce json
// @Param id path string true "ID"
// @Param request body product.UpdateProductRequest true "Produto"
// @Success 200 {object} product.ProductResponse
// @Security BearerAuth
// @Router /products/{id} [put]
func (h *ProductHandler) Update(
	w http.ResponseWriter,
	r *http.Request,
) {

	id := r.PathValue("id")

	var request productdto.UpdateProductRequest

	if err := httpx.ReadJSON(r, &request); err != nil {

		httpx.WriteError(
			w,
			http.StatusBadRequest,
			err,
		)

		return
	}

	product := mapper.ToProductUpdate(
		id,
		request,
	)

	if err := h.updateProductUseCase.Execute(product); err != nil {

		httpx.WriteError(
			w,
			http.StatusBadRequest,
			err,
		)

		return
	}

	httpx.WriteJSON(
		w,
		http.StatusOK,
		mapper.ToProductResponse(product),
	)
}

// Delete
//
// @Summary Remover Produto
// @Tags Products
// @Produce json
// @Param id path string true "ID"
// @Success 204
// @Security BearerAuth
// @Router /products/{id} [delete]
func (h *ProductHandler) Delete(
	w http.ResponseWriter,
	r *http.Request,
) {

	id := r.PathValue("id")

	if err := h.deleteProductUseCase.Execute(id); err != nil {

		httpx.WriteError(
			w,
			http.StatusBadRequest,
			err,
		)

		return
	}

	httpx.WriteStatus(
		w,
		http.StatusNoContent,
	)
}
