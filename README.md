
# Catalog API - Golang

Estrutura inicial em Golang utilizando Clean Architecture.

## Como rodar

1. Suba os containers:
```
docker-compose up --build
```

2. Acesse no navegador:
```
http://localhost:8080/health
```

Deve retornar:
```
Hello, World!
```

Banco de dados PostgreSQL disponível na porta `5432`.

## Estrutura de pastas

- `cmd/` → Main e inicialização
- `internal/domain/` → Entidades e interfaces
- `internal/usecase/` → Casos de uso
- `internal/infra/` → Infraestrutura (DB, etc)
- `internal/handler/` → HTTP handlers
- `docs/` → Documentação Swagger
- `migrations/` → Scripts SQL
# catalog
