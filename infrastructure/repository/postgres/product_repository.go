/*
struct ProductRepository

Responsabilidades:
- salvar produtos;
- atualizar produtos;
- remover produtos;
- buscar produto por ID;
- listar produtos.

Métodos:
- NewProductRepository()
- Save()
- Update()
- Delete()
- FindByID()
- FindAll()
*/

package postgres

import (
	"context"

	"desafio-go/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"

	"desafio-go/internal/domain"

	"desafio-go/infrastructure/database"

	"errors"

	"github.com/jackc/pgx/v5"
)

type ProductPostgresRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *database.Database) repository.ProductRepository {
	return &ProductPostgresRepository{
		db: db.Pool,
	}
}

func (r *ProductPostgresRepository) Save(product *domain.Product) error {

	query := `
		INSERT INTO products
			(id, name, price, stock)
		VALUES
			($1, $2, $3, $4)
	`

	_, err := r.db.Exec(
		context.Background(),
		query,
		product.ID,
		product.Name,
		product.Price,
		product.Stock,
	)

	return err
}

func (r *ProductPostgresRepository) Update(product *domain.Product) error {

	query := `
		UPDATE products
		SET
			name = $2,
			price = $3,
			stock = $4,
			updated_at = NOW()
		WHERE id = $1
	`

	_, err := r.db.Exec(
		context.Background(),
		query,
		product.ID,
		product.Name,
		product.Price,
		product.Stock,
	)

	return err
}

func (r *ProductPostgresRepository) Delete(id string) error {

	query := `
		DELETE FROM products
		WHERE id = $1
	`

	_, err := r.db.Exec(
		context.Background(),
		query,
		id,
	)

	return err
}

func (r *ProductPostgresRepository) FindByID(id string) (*domain.Product, error) {

	query := `
		SELECT
			id,
			name,
			price,
			stock
		FROM products
		WHERE id = $1
	`

	product := &domain.Product{}

	err := r.db.QueryRow(
		context.Background(),
		query,
		id,
	).Scan(
		&product.ID,
		&product.Name,
		&product.Price,
		&product.Stock,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrProductNotFound
	}

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (r *ProductPostgresRepository) FindAll() ([]*domain.Product, error) {

	query := `
		SELECT
			id,
			name,
			price,
			stock
		FROM products
		ORDER BY name
	`

	rows, err := r.db.Query(
		context.Background(),
		query,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*domain.Product

	for rows.Next() {

		product := &domain.Product{}

		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Price,
			&product.Stock,
		)
		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}
