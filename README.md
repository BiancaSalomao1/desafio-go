# Desafio Go API

API REST desenvolvida em Go utilizando **Clean Architecture**, princípios **SOLID**, autenticação **JWT**, **PostgreSQL**, **Docker** e documentação **Swagger**.

O objetivo do projeto é demonstrar boas práticas de desenvolvimento backend, separação de responsabilidades e implementação de regras de negócio em uma aplicação real.

---

# Funcionalidades

## Produtos

- Criar produto
- Buscar produto por ID
- Listar produtos
- Atualizar produto
- Remover produto

---

## Clientes

- Criar cliente
- Buscar cliente por ID
- Listar clientes
- Atualizar cliente
- Remover cliente

---

## Usuários

- Criar usuário
- Buscar usuário por ID
- Listar usuários
- Atualizar usuário
- Remover usuário
- Senhas armazenadas com BCrypt

---

## Pedidos

- Criar pedido
- Buscar pedido
- Listar pedidos
- Pagar pedido
- Cancelar pedido

---

## Autenticação

- Login
- JWT
- Middleware de autenticação
- Rotas protegidas

---

# Tecnologias

- Go 1.26
- PostgreSQL 17
- Docker
- Docker Compose
- JWT
- BCrypt
- Swagger (Swaggo)
- golang-migrate

---

# Arquitetura

O projeto segue o modelo **Clean Architecture**.

```
HTTP
        │
        ▼
Handlers
        │
        ▼
Use Cases
        │
        ▼
Repositories
        │
        ▼
PostgreSQL
```

Cada camada possui apenas uma responsabilidade e depende de abstrações, facilitando testes, manutenção e evolução do sistema.

---

# Estrutura do Projeto

```
cmd/
    api/

config/

docs/

infrastructure/
    database/
    http/
    repository/

internal/
    domain/
    dto/
    mapper/
    repository/
    security/
    usecase/

migrations/

scripts/

examples/
```

## cmd

Ponto de entrada da aplicação.

---

## config

Carregamento das configurações da aplicação.

---

## docs

Documentação Swagger gerada automaticamente.

---

## infrastructure

Implementações de infraestrutura.

- Banco de dados
- HTTP
- Repositórios PostgreSQL

---

## internal

Regra de negócio.

- Domain
- DTO
- Mapper
- Repository
- Security
- Use Cases

---

## migrations

Scripts SQL de criação do banco.

---

## scripts

Scripts auxiliares.

---

## examples

Exemplos de chamadas HTTP.

---

# Regras de Negócio

## Produto

- Nome obrigatório
- Preço maior que zero
- Estoque não negativo

---

## Cliente

- Nome obrigatório
- Email obrigatório

---

## Usuário

- Nome obrigatório
- Email obrigatório
- Senha obrigatória
- Senha armazenada utilizando BCrypt

---

## Pedido

- Cliente obrigatório
- Pelo menos um item
- Quantidade maior que zero
- Produto deve existir
- Estoque suficiente

Ao criar um pedido:

- estoque é reduzido

Ao cancelar:

- estoque retorna

Ao pagar:

- pedido torna-se imutável

---

# Executando Localmente

## Clonar

```bash
git clone <url>
```

```bash
cd desafio-go
```

---

## Configurar

Copie:

```bash
cp .env.example .env
```

---

## Instalar dependências

```bash
go mod tidy
```

---

## Executar

```bash
go run ./cmd/api
```

---

# Executando com Docker

```bash
docker compose up --build
```

A aplicação ficará disponível em:

```
http://localhost:8080
```

---

# Swagger

Após iniciar a aplicação:

```
http://localhost:8080/swagger/index.html
```

Caso a documentação seja alterada:

```bash
make swagger
```

---

# Makefile

Principais comandos:

```bash
make run

make build

make test

make vet

make fmt

make check

make swagger

make up

make down

make docker-build
```

---

# Scripts

Smoke Test

```bash
./scripts/smoke_test.sh
```

CRUD

```bash
./scripts/crud_test.sh
```

---

# Autenticação

Obter token:

```
POST /login
```

Utilizar:

```
Authorization: Bearer <TOKEN>
```

---

# Exemplos

Os exemplos completos encontram-se na pasta:

```
examples/
```

---

# Roadmap

## Concluído

- Clean Architecture
- SOLID
- PostgreSQL
- Docker
- JWT
- BCrypt
- Swagger
- Docker Compose
- Makefile
- GitHub Actions

---

## Próximas melhorias

- Testes unitários
- Cobertura de testes
- Health Check
- Version Endpoint
- Observabilidade
- Structured Logging

---

# Autor

Bianca Salomão

GitHub

https://github.com/BiancaSalomao1

LinkedIn

https://linkedin.com/in/bianca-salomao