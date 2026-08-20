
---

# `README.md`

Aqui eu também atualizaria. Principalmente porque hoje temos itens no **Roadmap como “Próximas melhorias” que já foram implementados**, como:

- testes unitários;
- Health Check;
- Structured Logging.

Minha sugestão de versão atual seria:

```md
# Desafio Go API

Projeto backend desenvolvido em Go com foco em arquitetura de software,
boas práticas e sistemas distribuídos.

A solução utiliza **Clean Architecture**, princípios **SOLID**, autenticação
com **JWT**, PostgreSQL, RabbitMQ e comunicação assíncrona entre serviços.

O projeto evoluiu de uma API REST para uma arquitetura distribuída, com
serviços independentes participando de fluxos de negócio através de mensageria
e Saga.

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
- Senhas protegidas com BCrypt

---

## Pedidos

- Criar pedido
- Buscar pedido por ID
- Listar pedidos
- Pagar pedido
- Cancelar pedido

---

## Autenticação

- Login
- Logout
- JWT
- Middleware de autenticação
- Rotas protegidas
- Controle de tokens ativos através de Token Store

---

# Arquitetura

O projeto segue os princípios da **Clean Architecture**.

```text
HTTP / Messaging
       │
       ▼
Handlers / Consumers
       │
       ▼
Use Cases
       │
       ▼
Domain
       │
       ▼
Repository Interfaces
       │
       ▼
Infrastructure