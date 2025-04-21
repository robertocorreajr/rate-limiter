# rate-limiter

Rate Limiter para requisições HTTP baseado em IP ou Token, com persistência em Redis.  
Ideal para controlar abuso de requisições por cliente ou chave de API.

## 🧠 Tecnologias utilizadas

- Go 1.24
- Redis 6.2
- Docker / Docker Compose
- Gorilla Mux
- Godotenv

## 🚀 Como executar

### 1. Clone o repositório

```bash
git clone https://github.com/robertocorreajr/rate-limiter.git
cd rate-limiter
```

### 2. Configure o ambiente

Crie um arquivo `.env` com as variáveis de ambiente:

```env
DEFAULT_LIMIT=10
BLOCK_TIME=300
REDIS_ADDR=redis:6379
```

### 3. Suba o ambiente com Docker

```bash
docker-compose up --build
```

O serviço estará acessível em:  
📍 `http://localhost:8080/`

## 🔌 Endpoints

### `GET /`

Requisição de teste protegida por rate limit.

#### Headers

- `API_KEY` (opcional): Token para controle individual

#### Respostas

- `200 OK`: Requisição aceita
- `429 Too Many Requests`: Limite excedido

## 🧪 Testes

### Testes automatizados (unitários)

Execute os testes unitários com:

```bash
go test ./internal/limiter
```

### Teste funcional com script

Você pode simular múltiplas requisições com:

```bash
go run scripts/test_rate_limit.go
```

> Esse script realiza chamadas com e sem token, testando limites e bloqueios.

## 🛠️ Desenvolvimento local (sem Docker)

```bash
go run cmd/rate-limiter/main.go
```

## 📁 Estrutura do projeto

```
rate-limiter/
├── cmd/                  # Ponto de entrada da aplicação
│   └── rate-limiter/
│       └── main.go
├── internal/             # Domínio e infraestrutura (privado ao módulo)
│   ├── handler/
│   ├── limiter/
│   └── storage/
├── scripts/              # Testes manuais e ferramentas
│   ├── test.go
│   └── test_rate_limit.go
├── Dockerfile
├── docker-compose.yml
├── go.mod / go.sum
├── README.md
└── .gitignore
```

## 🧠 Estratégia de persistência

A lógica de armazenamento do rate limit está desacoplada e implementada usando Redis.  
Você pode facilmente substituir o backend por outro adaptando a interface `Storage`.

## 📌 Notas finais

- Os testes utilizam o Redis local configurado no Docker Compose
- O projeto usa `.env` para facilitar configuração e evitar hardcode
- Está pronto para escalar, com arquitetura organizada e extensível

---

Feito com 💻 por [Roberto Corrêa](https://github.com/robertocorreajr)
