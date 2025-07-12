
# 📦 Catalog API

Uma API desenvolvida em **Go** com foco em **Clean Architecture**, utilizando **PostgreSQL**, **Redis** para cache e um sistema de **eventos customizado**. Essa API tem como objetivo gerenciar **nomes de combos**, podendo ser estendida para representar um catálogo completo de produtos.

---

## ✨ Funcionalidades

- 🔎 Listar combos (com paginação)
- 📄 Obter combo por ID
- ➕ Criar novo combo
- ✏️ Atualizar combo
- 🚫 Desabilitar combo
- 🚀 Disparar eventos no momento da criação (`ComboNameCreated`)
- ⚡ Cache Redis para operações (em construção)

---

## 📐 Arquitetura

O projeto segue os princípios da **Clean Architecture**:

```
internal/
├── entity/          # Entidades do domínio
├── usecase/         # Casos de uso (regras de negócio)
├── infra/           # Infraestrutura (DB, cache, handlers HTTP, etc)
│   ├── database/
│   ├── cache/
│   ├── repository/
│   └── web/
│       ├── handler/
│       └── server/
├── event/           # Manipuladores de eventos
pkg/
└── events/          # Definição de interfaces de eventos
```

---

## 🧪 Testes

Testes automatizados escritos com `testify`.

```bash
go test ./... -v
```

Geração de cobertura:

```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

---

## 🛠️ Tecnologias

- [Golang](https://golang.org)
- [PostgreSQL](https://www.postgresql.org/)
- [Redis](https://redis.io/)
- [Testify](https://github.com/stretchr/testify)
- Arquitetura Limpa (Clean Architecture)

---

## ⚙️ Como rodar localmente

### 1. Clone o projeto

```bash
git clone https://github.com/seu-usuario/catalog.git
cd catalog
```

### 2. Suba os serviços com Docker Compose

```bash
docker-compose up -d
```

Isso iniciará:
- PostgreSQL
- Redis
- A aplicação `catalog`

### 3. Acesse a API

Por padrão, estará rodando em:

```
http://localhost:8080
```

---

## 🧪 Endpoints

| Método | Rota                     | Descrição                  |
|--------|--------------------------|----------------------------|
| GET    | /combonames              | Lista todos os combos      |
| GET    | /combonames/{id}         | Retorna combo por ID       |
| POST   | /combonames              | Cria um novo combo         |
| PUT    | /combonames/{id}         | Atualiza um combo          |
| DELETE | /combonames/{id}         | Desabilita um combo        |

> A documentação Swagger pode ser adicionada futuramente.

---

## 📚 Futuras melhorias

- ✅ Suporte completo a Redis (read-through + invalidation)
- ✅ Geração automática de UUID
- 📖 Swagger/OpenAPI
- 🔐 Autenticação via JWT
- 🌐 gRPC ou GraphQL interface

---

## 🧑‍💻 Autor

**Lucas Batista de Souza**  
Tech Manager | Backend Developer  
📧 lucas_araujo91@hotmail.com

---

## 🪪 Licença

MIT
