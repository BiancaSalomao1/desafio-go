# Arquitetura

Este projeto evoluiu de uma API REST monolítica para uma solução distribuída,
mantendo os princípios de **Clean Architecture**, **SOLID** e separação clara
entre domínio, aplicação e infraestrutura.

A solução é composta por serviços independentes que se comunicam
assincronamente através do RabbitMQ em determinados fluxos de negócio.

---

# Clean Architecture

Cada serviço mantém suas regras de negócio independentes de detalhes externos
como HTTP, banco de dados, mensageria ou cache.

```mermaid
flowchart TB

    Client["Cliente HTTP"]

    Router["Router"]

    Middleware["Middlewares"]

    Handler["HTTP Handler"]

    DTO["DTO / Mapper"]

    UseCase["Use Case"]

    Repository["Repository Interface"]

    Infrastructure["Infrastructure"]

    Database[("PostgreSQL")]

    Client --> Router
    Router --> Middleware
    Middleware --> Handler
    Handler --> DTO
    DTO --> UseCase
    UseCase --> Repository
    Repository --> Infrastructure
    Infrastructure --> Database
```

---

# Fluxo de Autenticação

O sistema utiliza JWT para autenticação.

```mermaid
sequenceDiagram

    actor Client

    participant API
    participant UserRepository
    participant BCrypt
    participant JWT
    participant TokenStore

    Client->>API: POST /login

    API->>UserRepository: FindByEmail()
    UserRepository-->>API: User

    API->>BCrypt: CheckPassword()
    BCrypt-->>API: Password valid

    API->>JWT: GenerateToken()
    JWT-->>API: Token

    API->>TokenStore: Save(token, TTL)
    TokenStore-->>API: Token active

    API-->>Client: accessToken
```

---

# Middleware de Autenticação

As rotas protegidas passam pelo middleware de autenticação.

O middleware valida:

- Presença do header `Authorization`;
- Formato `Bearer <token>`;
- Validade do JWT;
- Existência do token ativo no `TokenStore`.

```mermaid
sequenceDiagram

    actor Client

    participant Middleware
    participant JWT
    participant TokenStore
    participant Handler

    Client->>Middleware: Request + Bearer Token

    Middleware->>JWT: ValidateToken()
    JWT-->>Middleware: Claims

    Middleware->>TokenStore: Exists(token)
    TokenStore-->>Middleware: Token active

    Middleware->>Handler: Request com Context
```

---

# Fluxo de Logout

O logout invalida o token ativo.

Após o logout, o token não deve mais ser considerado válido.

```mermaid
sequenceDiagram

    actor Client

    participant Handler
    participant LogoutUseCase
    participant TokenStore

    Client->>Handler: POST /logout + Bearer Token

    Handler->>LogoutUseCase: Execute(ctx, token)

    LogoutUseCase->>TokenStore: Delete(ctx, token)

    TokenStore-->>LogoutUseCase: OK

    LogoutUseCase-->>Handler: Success

    Handler-->>Client: 204 No Content
```

---

# Fluxo de Criação de Pedido

A criação de pedidos envolve validações de domínio e persistência.

Quando o fluxo exige comunicação com outro serviço, a integração ocorre
através do broker.

```mermaid
flowchart LR

    Client["Cliente"]

    Handler["Order Handler"]

    UseCase["CreateOrder Use Case"]

    CustomerRepository["Customer Repository"]

    OrderRepository["Order Repository"]

    Publisher["Event Publisher"]

    Database[("Orders Database")]

    Client --> Handler

    Handler --> UseCase

    UseCase --> CustomerRepository
    UseCase --> OrderRepository
    UseCase --> Publisher

    CustomerRepository --> Database
    OrderRepository --> Database
```

---

# Comunicação Assíncrona

A comunicação entre os serviços utiliza RabbitMQ.

Os serviços não acessam diretamente o banco de dados uns dos outros.

Cada serviço é responsável pelos seus próprios dados.

```mermaid
flowchart LR

    OrdersAPI["Orders API"]

    RabbitMQ["RabbitMQ"]

    ProductService["Product Service"]

    OrdersAPI -->|Publish Event| RabbitMQ

    RabbitMQ -->|Consume Event| ProductService

    ProductService -->|Publish Result| RabbitMQ

    RabbitMQ -->|Consume Result| OrdersAPI
```

---

# Saga

O fluxo distribuído é implementado utilizando uma Saga.

A Saga permite manter consistência entre serviços sem utilizar transações
distribuídas.

```mermaid
sequenceDiagram

    participant Orders as Orders API
    participant Broker as RabbitMQ
    participant Products as Product Service

    Orders->>Broker: OrderCreated

    Broker->>Products: Process Order

    alt Estoque disponível

        Products->>Broker: StockReserved

        Broker->>Orders: StockReserved

        Orders->>Orders: Confirm Order

    else Falha no estoque

        Products->>Broker: StockReservationFailed

        Broker->>Orders: StockReservationFailed

        Orders->>Orders: Compensate Order

    end
```

---

# Idempotência

Mensagens podem ser entregues mais de uma vez pelo broker.

Por esse motivo, o processamento deve impedir que uma mensagem repetida
corrompa o estado da aplicação.

```text
Message ID
     │
     ▼
Já processada?
     │
 ┌───┴────┐
 │        │
Sim      Não
 │        │
 ▼        ▼
Ignore   Process
```

---

# Estrutura da Solução

```text
orders-api/

├── cmd/
│   └── api/
│
├── config/
│
├── docs/
│
├── infrastructure/
│   ├── cache/
│   ├── database/
│   ├── http/
│   │   ├── handler/
│   │   ├── middleware/
│   │   └── routes/
│   ├── logging/
│   ├── messaging/
│   │   └── rabbitmq/
│   └── repository/
│       ├── memory/
│       └── postgres/
│
├── internal/
│   ├── app/
│   ├── domain/
│   ├── dto/
│   ├── mapper/
│   ├── messaging/
│   ├── repository/
│   ├── security/
│   └── usecase/
│
├── migrations/
│
├── scripts/
│
├── examples/
│
└── tests/
    └── integration/
```

---

# Princípios Aplicados

## SOLID

### Single Responsibility Principle

Cada componente possui uma responsabilidade principal.

Exemplos:

- Handler → HTTP;
- Use Case → regra de aplicação;
- Repository → persistência;
- Entity → regras do domínio;
- Middleware → processamento transversal das requisições;
- TokenStore → gerenciamento do estado dos tokens.

---

### Open/Closed Principle

Novas implementações podem ser adicionadas sem modificar os casos de uso.

Exemplo:

```text
UserRepository
      │
      ├── MemoryUserRepository
      │
      └── PostgreSQLUserRepository
```

---

### Liskov Substitution Principle

Qualquer implementação compatível com uma interface pode substituir outra.

Exemplo:

```text
TokenStore
    │
    ├── MemoryTokenStore
    │
    └── RedisTokenStore
```

---

### Interface Segregation Principle

As interfaces representam contratos específicos da aplicação.

Exemplos:

- `UserRepository`;
- `OrderRepository`;
- `ProductRepository`;
- `TokenStore`;
- interfaces específicas dos Use Cases utilizadas pelos handlers.

---

### Dependency Inversion Principle

Os casos de uso dependem de abstrações, e não de implementações concretas.

```mermaid
flowchart TB

    UseCase["Use Case"]

    RepositoryInterface["Repository Interface"]

    PostgreSQL["PostgreSQL Repository"]

    Memory["Memory Repository"]

    UseCase --> RepositoryInterface

    RepositoryInterface --> PostgreSQL
    RepositoryInterface --> Memory
```