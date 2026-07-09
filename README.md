# Desafio Go - Serviço de Pedidos

API REST desenvolvida em **Go** para gerenciamento de pedidos de uma loja.

O projeto foi desenvolvido de forma incremental durante a disciplina, evoluindo desde um domínio executado em memória até uma API REST completa utilizando PostgreSQL, Docker e boas práticas de arquitetura.

---

# Objetivos

Desenvolver uma API REST capaz de:

- cadastrar usuários;
- autenticar usuários (JWT opcional);
- cadastrar clientes;
- listar clientes;
- buscar clientes por ID;
- cadastrar produtos;
- listar produtos;
- buscar produtos por ID;
- criar pedidos;
- listar pedidos;
- consultar pedidos;
- pagar pedidos;
- cancelar pedidos.

---

# Evolução do Projeto

O desenvolvimento será dividido em três etapas.

## Aula 03 – Domínio

Implementação apenas da lógica de negócio utilizando a biblioteca padrão do Go.

### Características

- domínio
- interfaces
- repositories em memória
- services (use cases)
- tratamento de erros

Sem:

- HTTP
- JSON
- banco de dados
- Docker
- frameworks

---

## Aula 06 – API REST

Evolução do domínio para uma API REST.

Novos recursos:

- PostgreSQL
- Docker
- pgx
- migrations
- handlers
- DTOs
- paginação
- transações

---

## Aula 08 – Projeto Integrador

Refinamento da aplicação.

Novos recursos:

- testes unitários
- logs
- context
- concorrência
- documentação Swagger
- melhorias arquiteturais

---

# Arquitetura

O projeto segue uma arquitetura inspirada em **Clean Architecture** e **Ports & Adapters**.

```text
HTTP Request
        │
        ▼
 Routes (Chi)
        │
        ▼
 Middleware
        │
        ▼
 Handlers
        │
        ▼
 DTO
        │
        ▼
 Use Cases
        │
        ▼
 Repository Interfaces
        │
        ▼
 PostgreSQL Repositories
        │
        ▼
 pgxpool
        │
        ▼
 PostgreSQL
```

As regras de negócio permanecem independentes da camada HTTP e da infraestrutura.

---

# Estrutura do Projeto

```text
pedidos/

├── cmd/
│   └── api/
│       └── main.go
│
├── config/
│
├── database/
│
├── docs/
│
├── internal/
│
│   ├── domain/
│   │   ├── user.go
│   │   ├── customer.go
│   │   ├── product.go
│   │   ├── order.go
│   │   ├── order_item.go
│   │   └── errors.go
│   │
│   ├── dto/
│   │   ├── user/
│   │   ├── customer/
│   │   ├── product/
│   │   └── order/
│   │
│   ├── handler/
│   │   ├── user_handler.go
│   │   ├── customer_handler.go
│   │   ├── product_handler.go
│   │   └── order_handler.go
│   │
│   ├── middleware/
│   │
│   ├── repository/
│   │   ├── interfaces/
│   │   └── postgres/
│   │
│   ├── routes/
│   │
│   └── usecase/
│
├── migrations/
│
├── Makefile
│
├── docker-compose.yml
│
├── go.mod
│
└── go.sum
```

---

# Funcionalidades

## Usuários

- cadastrar usuário;
- listar usuários;
- buscar usuário por ID;
- autenticar usuário (JWT opcional).

---

## Clientes

- cadastrar cliente;
- listar clientes;
- buscar cliente por ID.

---

## Produtos

- cadastrar produto;
- listar produtos;
- buscar produto por ID.

---

## Pedidos

- criar pedido;
- listar pedidos;
- buscar pedido;
- pagar pedido;
- cancelar pedido.

---

# Regras de Negócio

## Usuário

- nome obrigatório;
- e-mail obrigatório;
- e-mail único;
- senha armazenada apenas como `password_hash`;
- autenticação utilizando JWT (opcional).

---

## Cliente

- nome obrigatório;
- e-mail obrigatório.

---

## Produto

- nome obrigatório;
- preço maior que zero;
- estoque maior ou igual a zero.

---

## Pedido

- cliente obrigatório;
- cliente deve existir;
- pedido deve possuir pelo menos um item;
- quantidade maior que zero;
- produto deve existir;
- estoque suficiente;
- preço do produto congelado no momento da compra;
- redução automática do estoque;
- status inicial `PENDING`;
- pagamento altera para `PAID`;
- cancelamento altera para `CANCELED`;
- cancelamento devolve o estoque;
- pedidos pagos ou cancelados não podem alterar novamente o status.

---

# Modelo Relacional

```text
Users

Customers
      │
      │ 1
      │
      └───────────────┐
                      │ N
                  Orders
                      │
                      │ 1
                      │
                      │ N
                OrderItems
                      │
                      │ N
                      │
                      │ 1
                  Products
```

---

# Banco de Dados

Tabelas previstas:

- users
- customers
- products
- orders
- order_items

---

# Docker

Configuração padrão utilizada na disciplina.

```yaml
services:

  postgres:

    image: postgres:18.4

    container_name: postgres18

    environment:
      POSTGRES_USER: app
      POSTGRES_PASSWORD: app
      POSTGRES_DB: app

    ports:
      - "${POSTGRES_PORT:-5432}:5432"

    volumes:
      - postgres_data:/var/lib/postgresql

volumes:

  postgres_data:
```

> **Observação**
>
> A disciplina utiliza a porta **5432** como padrão.
>
> No meu ambiente de desenvolvimento (Linux), essa porta já é utilizada por uma instalação local do PostgreSQL.
>
> Portanto será utilizado um arquivo `.env` contendo:
>
> ```env
> POSTGRES_PORT=5433
> ```
>
> Dessa forma, basta remover essa variável para utilizar a porta padrão da disciplina.

---

# Migrations

As tabelas serão criadas utilizando **golang-migrate**.

```text
migrations/

000001_create_users

000002_create_customers

000003_create_products

000004_create_orders

000005_create_order_items
```

---

# Endpoints

## Usuários

```http
POST   /usuarios
GET    /usuarios
GET    /usuarios/{id}
POST   /login
```

---

## Clientes

```http
POST   /clientes
GET    /clientes
GET    /clientes/{id}
```

---

## Produtos

```http
POST   /produtos
GET    /produtos
GET    /produtos/{id}
```

---

## Pedidos

```http
POST   /pedidos
GET    /pedidos
GET    /pedidos/{id}
POST   /pedidos/{id}/pagar
POST   /pedidos/{id}/cancelar
```

---

# Segurança

A autenticação utilizando **JWT** é considerada **opcional** neste projeto.

Caso implementada, será utilizada para autenticação de usuários sem alterar a arquitetura principal.

---

# Documentação da API

Será utilizada a biblioteca **Swaggo** para geração automática da documentação Swagger/OpenAPI.

Após configurada:

```text
http://localhost:8080/swagger/index.html
```

---

# Roteamento HTTP

Será utilizado o **Chi Router** para gerenciamento das rotas HTTP e middlewares.

---

# Testes

Serão implementados testes unitários para as principais regras de negócio.

Exemplos:

- criação de pedido;
- pagamento;
- cancelamento;
- estoque insuficiente;
- cliente inexistente;
- produto inexistente.

---

# Dependências

```bash
go get github.com/jackc/pgx/v5

go get github.com/go-chi/chi/v5

go get github.com/golang-migrate/migrate/v4

go get github.com/google/uuid

go get github.com/swaggo/swag

go get golang.org/x/crypto/bcrypt
```

Caso algum import apresente erro:

```bash
go mod tidy
```

Para utilizar o comando `migrate` pela linha de comando:

```bash
go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

---

# Makefile

Os principais comandos do projeto serão automatizados.

```bash
make up

make down

make run

make migrate-up

make migrate-down

make test

make fmt

make tidy

make vet
```

---

# Convenção de Commits

O projeto seguirá o padrão **Conventional Commits**.

Exemplos:

```bash
git commit -m "docs: adiciona README inicial"

git commit -m "feat: cria estrutura inicial do projeto"

git commit -m "feat: implementa entidade Product"

git commit -m "feat: implementa repository em memória"

git commit -m "feat: implementa criação de pedidos"

git commit -m "feat: adiciona API REST"

git commit -m "feat: integra PostgreSQL"

git commit -m "feat: implementa migrations"

git commit -m "feat: implementa autenticação JWT"

git commit -m "docs: documenta endpoints Swagger"

git commit -m "test: adiciona testes unitários"

git commit -m "fix: corrige validação de estoque"

git commit -m "refactor: reorganiza arquitetura"

git commit -m "style: aplica gofmt"
```

---

# Referências

- **Conventional Commits**  
  https://www.conventionalcommits.org/en/v1.0.0/

- **Go — Documentação Oficial**  
  https://go.dev/doc/

- **Effective Go**  
  https://go.dev/doc/effective_go

- **The Twelve-Factor App**  
  https://12factor.net/pt_br/

- **PostgreSQL**  
  https://www.postgresql.org/docs/current/

- **UUID no PostgreSQL**  
  https://www.postgresql.org/docs/current/functions-uuid.html

- **pgx**  
  https://github.com/jackc/pgx

- **golang-migrate**  
  https://github.com/golang-migrate/migrate

- **Chi Router**  
  https://github.com/go-chi/chi

- **Swaggo**  
  https://github.com/swaggo/swag

- **Docker**  
  https://www.docker.com/

- **Docker Compose**  
  https://docs.docker.com/compose/