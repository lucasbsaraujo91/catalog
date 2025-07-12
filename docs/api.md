
# 📘 Documentação da API - Catálogo

Todas as rotas (exceto /ping e /login) exigem autenticação via token fixo.

## 🔐 Autenticação

Adicione o seguinte header a todas as requisições:

```
Authorization: Bearer meu-token-super-secreto
```

---

## 📡 Endpoints

### ✅ `GET /ping`
- **Descrição:** Teste simples para verificar se o serviço está no ar
- **Auth:** ❌ Não requer autenticação
- **Resposta:**
```text
pong
```

---

### 🔐 `POST /login`
- **Descrição:** Login via usuário fixo (admin / 123456)
- **Auth:** ❌ Não requer autenticação
- **Corpo da requisição (JSON):**
```json
{
  "username": "admin",
  "password": "123456"
}
```
- **Resposta:**
```json
{
  "token": "meu-token-super-secreto"
}
```

---

### 🔐 `GET /combo-names`
- **Descrição:** Lista todos os ComboNames disponíveis
- **Auth:** ✅ Requer token fixo
- **Resposta:**
```json
[
  {
    "id": 1,
    "name": "Combo A",
    "active": true
  },
  ...
]
```

---

### 🔐 `GET /combo-names/{id}`
- **Descrição:** Busca um combo pelo ID
- **Auth:** ✅
- **Resposta de sucesso:**
```json
{
  "id": 1,
  "name": "Combo A",
  "active": true
}
```
- **Erro 404:** se não encontrado

---

### 🔐 `POST /combo-names`
- **Descrição:** Cria um novo combo
- **Auth:** ✅
- **Corpo da requisição:**
```json
{
  "name": "Combo Novo"
}
```
- **Resposta 201:**
```json
{
  "id": 123,
  "name": "Combo Novo",
  "active": true
}
```

---

### 🔐 `PUT /combo-names/{id}`
- **Descrição:** Atualiza um combo existente
- **Auth:** ✅
- **Corpo esperado:**
```json
{
  "name": "Combo Atualizado"
}
```

---

### 🔐 `DELETE /combo-names/{id}`
- **Descrição:** Desativa um combo
- **Auth:** ✅

---

### 🔐 `POST /limpa-cache`
- **Descrição:** Limpa manualmente o cache do Redis
- **Auth:** ✅
- **Resposta:**
```json
{
  "message": "Cache limpo com sucesso"
}
```

---

## ⚠️ Códigos de resposta comuns

| Código | Significado         |
|--------|---------------------|
| 200    | OK                  |
| 201    | Criado              |
| 400    | Requisição inválida |
| 401    | Não autorizado      |
| 404    | Não encontrado      |
| 500    | Erro interno        |
