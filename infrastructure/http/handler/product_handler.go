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

	"desafio-go/infrastructure/http/httpx"

	productdto "desafio-go/internal/dto/product"
	"desafio-go/internal/mapper"

	productusecase "desafio-go/internal/usecase/product"
)

type ProductHandler struct {
	createProductUseCase *productusecase.CreateProductUseCase
	getProductUseCase    *productusecase.GetProductUseCase
	listProductsUseCase  *productusecase.ListProductsUseCase
	updateProductUseCase *productusecase.UpdateProductUseCase
	deleteProductUseCase *productusecase.DeleteProductUseCase
}

func NewProductHandler(
	createProductUseCase *productusecase.CreateProductUseCase,
	getProductUseCase *productusecase.GetProductUseCase,
	listProductsUseCase *productusecase.ListProductsUseCase,
	updateProductUseCase *productusecase.UpdateProductUseCase,
	deleteProductUseCase *productusecase.DeleteProductUseCase,
) *ProductHandler {

	return &ProductHandler{
		createProductUseCase: createProductUseCase,
		getProductUseCase:    getProductUseCase,
		listProductsUseCase:  listProductsUseCase,
		updateProductUseCase: updateProductUseCase,
		deleteProductUseCase: deleteProductUseCase,
	}
}

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

/*
Função Delete

Responsabilidades:
- remover um produto.

Métodos:
- Delete()
*/

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
