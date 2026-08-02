package postgres

/*
struct CustomerPostgresRepository

Responsabilidades:
- salvar cliente;
- atualizar cliente;
- remover cliente;
- buscar cliente por ID;
- listar clientes.

Métodos:
- NewCustomerRepository()
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

	"github.com/jackc/pgx/v5"
)

type CustomerPostgresRepository struct {
	db repository.DBTX
}

func NewCustomerRepository(db repository.DBTX) repository.CustomerRepository {
	return &CustomerPostgresRepository{
		db: db,
	}
}

func (r *CustomerPostgresRepository) Save(customer *domain.Customer) error {

	query := `
INSERT INTO customers
(
	id,
	name,
	email,
	password_hash
)
VALUES
(
	$1,
	$2,
	$3,
	$4
)
`

	_, err := r.db.Exec(
		context.Background(),
		query,
		customer.ID,
		customer.Name,
		customer.Email,
		customer.PasswordHash,
	)

	return err
}

func (r *CustomerPostgresRepository) Update(customer *domain.Customer) error {

	query := `
		UPDATE customers
		SET
			name = $2,
			email = $3,
			password_hash = $4,
			updated_at = NOW()
		WHERE id = $1
	`

	_, err := r.db.Exec(
		context.Background(),
		query,
		customer.ID,
		customer.Name,
		customer.Email,
		customer.PasswordHash,
	)

	return err
}

func (r *CustomerPostgresRepository) Delete(id string) error {

	_, err := r.db.Exec(
		context.Background(),
		`DELETE FROM customers WHERE id = $1`,
		id,
	)

	return err
}

func (r *CustomerPostgresRepository) FindByID(id string) (*domain.Customer, error) {

	customer := &domain.Customer{}

	err := r.db.QueryRow(
		context.Background(),
		`SELECT id,name,email, password_hash
		 FROM customers
		 WHERE id=$1`,
		id,
	).Scan(
		&customer.ID,
		&customer.Name,
		&customer.Email,
		&customer.PasswordHash,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrCustomerNotFound
	}

	if err != nil {
		return nil, err
	}

	return customer, nil
}

func (r *CustomerPostgresRepository) FindAll() ([]*domain.Customer, error) {

	rows, err := r.db.Query(
		context.Background(),
		`SELECT id,name,email,password_hash
		 FROM customers
		 ORDER BY name`,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var customers []*domain.Customer

	for rows.Next() {

		customer := &domain.Customer{}

		err := rows.Scan(
			&customer.ID,
			&customer.Name,
			&customer.Email,
			&customer.PasswordHash,
		)

		if err != nil {
			return nil, err
		}

		customers = append(customers, customer)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return customers, nil
}

func (r *CustomerPostgresRepository) FindByEmail(email string) (*domain.Customer, error) {

	customer := &domain.Customer{}

	err := r.db.QueryRow(
		context.Background(),
		`SELECT id,name,email,password_hash
		 FROM customers
		 WHERE email=$1`,
		email,
	).Scan(
		&customer.ID,
		&customer.Name,
		&customer.Email,
		&customer.PasswordHash,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, domain.ErrCustomerNotFound
	}

	if err != nil {
		return nil, err
	}

	return customer, nil
}
