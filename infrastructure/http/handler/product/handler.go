package handler

/*
struct ProductHandler

Responsabilidades:
- receber requisições HTTP relacionadas a produtos;
- converter DTOs para entidades de domínio;
- executar os casos de uso;
- retornar respostas HTTP em JSON.

Campos:
- createProductUseCase
- getProductUseCase
- listProductsUseCase

Métodos:
- NewProductHandler()
- Create()
- GetByID()
- List()
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
}

func NewProductHandler(
	createProductUseCase *productusecase.CreateProductUseCase,
	getProductUseCase *productusecase.GetProductUseCase,
	listProductsUseCase *productusecase.ListProductsUseCase,
) *ProductHandler {

	return &ProductHandler{
		createProductUseCase: createProductUseCase,
		getProductUseCase:    getProductUseCase,
		listProductsUseCase:  listProductsUseCase,
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
