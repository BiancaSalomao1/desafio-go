package postgres

/*
struct OrderPostgresRepository

Responsabilidades:
- salvar pedido;
- atualizar pedido;
- remover pedido;
- buscar pedido por ID;
- listar pedidos.

Métodos:
- NewOrderRepository()
- Save()
- Update()
- Delete()
- FindByID()
- FindAll()
*/

import (
	"context"
	"errors"

	"desafio-go/internal/domain"
	"desafio-go/internal/repository"

	"github.com/google/uuid"
)

type OrderPostgresRepository struct {
	db repository.DBTX
}

func NewOrderRepository(db repository.DBTX) repository.OrderRepository {
	return &OrderPostgresRepository{
		db: db,
	}
}

func (r *OrderPostgresRepository) Save(order *domain.Order) error {

	ctx := context.Background()

	_, err := r.db.Exec(
		ctx,
		`
		INSERT INTO orders
		(id, customer_id, status)
		VALUES ($1,$2,$3)
		`,
		order.ID,
		order.CustomerID,
		order.Status,
	)

	if err != nil {
		return err
	}

	for _, item := range order.Items {

		_, err = r.db.Exec(
			ctx,
			`
			INSERT INTO order_items
			(
				id,
				order_id,
				product_id,
				product_name,
				product_price,
				quantity
			)
			VALUES
			($1,$2,$3,$4,$5,$6)
			`,
			uuid.NewString(),
			order.ID,
			item.ProductID,
			item.Name,
			item.Price,
			item.Quantity,
		)

		if err != nil {
			return err
		}
	}

	return nil
}

func (r *OrderPostgresRepository) Update(order *domain.Order) error {

	_, err := r.db.Exec(
		context.Background(),
		`
		UPDATE orders
		SET
			status = $2,
			updated_at = NOW()
		WHERE id = $1
		`,
		order.ID,
		order.Status,
	)

	return err
}

func (r *OrderPostgresRepository) Delete(id string) error {

	_, err := r.db.Exec(
		context.Background(),
		`DELETE FROM orders WHERE id = $1`,
		id,
	)

	return err
}

func (r *OrderPostgresRepository) FindByID(id string) (*domain.Order, error) {
	return nil, errors.New("not implemented")
}

func (r *OrderPostgresRepository) FindAll() ([]*domain.Order, error) {
	return nil, errors.New("not implemented")
}
