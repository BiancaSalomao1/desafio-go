package integration

/*
database_test

Responsabilidades:

- abrir conexão com o PostgreSQL;
- limpar todas as tabelas;
- garantir isolamento entre os testes;
- disponibilizar funções auxiliares para acesso ao banco.

Este arquivo NÃO testa endpoints.

Seu objetivo é garantir que cada teste execute em um banco limpo.
*/

import (
	"context"
	"desafio-go/infrastructure/database"
	"testing"

	"github.com/jackc/pgx/v5"
)

const connectionString = "postgres://postgres:postgres@localhost:5432/desafio_go"

func connectDatabase(
	t *testing.T,
) *pgx.Conn {

	t.Helper()

	db, err := pgx.Connect(
		context.Background(),
		connectionString,
	)

	if err != nil {
		t.Fatalf("error connecting database: %v", err)
	}

	return db
}

func cleanDatabase(
	t *testing.T,
	db *database.Database,
) {
	t.Helper()

	_, err := db.Pool.Exec(
		context.Background(),
		`
		TRUNCATE TABLE
			order_items,
			orders,
			products,
			customers,
			users
		RESTART IDENTITY CASCADE;
		`,
	)

	if err != nil {
		t.Fatalf("error cleaning database: %v", err)
	}
}
