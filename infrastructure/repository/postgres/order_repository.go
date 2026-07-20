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
	"github.com/jackc/pgx/v5"
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

	for i := range order.Items {

		itemID := uuid.NewString()

		order.Items[i].ID = itemID

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
			itemID,
			order.ID,
			order.Items[i].ProductID,
			order.Items[i].Name,
			order.Items[i].Price,
			order.Items[i].Quantity,
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

	ctx := context.Background()

	order := &domain.Order{}

	query := `
		SELECT
			id,
			customer_id,
			status
		FROM orders
		WHERE id = $1
	`

	err := r.db.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&order.ID,
		&order.CustomerID,
		&order.Status,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrOrderNotFound
	}

	if err != nil {
		return nil, err
	}

	queryItems := `
		SELECT
			id,
			product_id,
			product_name,
			product_price,
			quantity
		FROM order_items
		WHERE order_id = $1
		ORDER BY product_name
	`

	rows, err := r.db.Query(
		ctx,
		queryItems,
		order.ID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {

		item := domain.OrderItem{}

		err := rows.Scan(
			&item.ID,
			&item.ProductID,
			&item.Name,
			&item.Price,
			&item.Quantity,
		)

		if err != nil {
			return nil, err
		}

		order.Items = append(order.Items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return order, nil
}

func (r *OrderPostgresRepository) FindAll() ([]*domain.Order, error) {

	ctx := context.Background()

	query := `
		SELECT
			id,
			customer_id,
			status
		FROM orders
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*domain.Order

	for rows.Next() {

		order := &domain.Order{}

		err := rows.Scan(
			&order.ID,
			&order.CustomerID,
			&order.Status,
		)
		if err != nil {
			return nil, err
		}

		itemQuery := `
			SELECT
				id,
				product_id,
				product_name,
				product_price,
				quantity
			FROM order_items
			WHERE order_id = $1
			ORDER BY product_name
		`

		itemRows, err := r.db.Query(
			ctx,
			itemQuery,
			order.ID,
		)
		if err != nil {
			return nil, err
		}

		for itemRows.Next() {

			item := domain.OrderItem{}

			err := itemRows.Scan(
				&item.ID,
				&item.ProductID,
				&item.Name,
				&item.Price,
				&item.Quantity,
			)
			if err != nil {
				itemRows.Close()
				return nil, err
			}

			order.Items = append(order.Items, item)
		}

		if err := itemRows.Err(); err != nil {
			itemRows.Close()
			return nil, err
		}

		itemRows.Close()

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *OrderPostgresRepository) List(limit, offset int) ([]*domain.Order, error) {

	ctx := context.Background()

	query := `
		SELECT
			id,
			customer_id,
			status
		FROM orders
		ORDER BY created_at DESC
		LIMIT $1
		OFFSET $2
	`

	rows, err := r.db.Query(
		ctx,
		query,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*domain.Order

	for rows.Next() {

		order := &domain.Order{}

		err := rows.Scan(
			&order.ID,
			&order.CustomerID,
			&order.Status,
		)
		if err != nil {
			return nil, err
		}

		itemQuery := `
			SELECT
				id,
				product_id,
				product_name,
				product_price,
				quantity
			FROM order_items
			WHERE order_id = $1
			ORDER BY product_name
		`

		itemRows, err := r.db.Query(
			ctx,
			itemQuery,
			order.ID,
		)
		if err != nil {
			return nil, err
		}

		for itemRows.Next() {

			item := domain.OrderItem{}

			err := itemRows.Scan(
				&item.ID,
				&item.ProductID,
				&item.Name,
				&item.Price,
				&item.Quantity,
			)
			if err != nil {
				itemRows.Close()
				return nil, err
			}

			order.Items = append(order.Items, item)
		}

		if err := itemRows.Err(); err != nil {
			itemRows.Close()
			return nil, err
		}

		itemRows.Close()

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}
