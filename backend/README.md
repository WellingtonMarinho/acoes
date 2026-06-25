# Backend

API em Go para alertas de ações da B3.

## Como subir localmente

Na raiz do projeto:

```bash
make run-backend
```

Isso sobe o backend em `http://localhost:8080` e um Postgres em `localhost:5432`.
Quando `DATABASE_URL` estiver definido, o backend usa o Postgres para alertas e devices.
No fluxo com Docker Compose, o serviço de migrations aplica o schema antes da API subir.
As migrations seguem o padrão do `goose` e ficam em `internal/postgres/migrations/`.

Para desenvolvimento com reload automatico do backend ao salvar arquivos `.go`, use:

```bash
make dev-backend
```

Esse modo monta `backend/` dentro do container, aplica as migrations ao iniciar e usa `air` para recompilar/reiniciar o servidor sem rebuildar a imagem a cada alteração de código.
Se `Dockerfile.dev` mudar, reconstrua a imagem dev com `make build-dev-backend`.

## Variáveis de ambiente

- `HTTP_ADDR`: endereço HTTP do servidor. Padrão: `:8080`
- `MONITOR_INTERVAL_SECONDS`: intervalo do worker de monitoramento. Padrão: `10`
- `ALERTS_STORE_PATH`: caminho do arquivo de persistência dos alertas
- `DEVICES_STORE_PATH`: caminho do arquivo de persistência dos devices
- `DATABASE_URL`: string de conexão do Postgres. Quando definida, o backend usa persistência relacional
- `PRICEFEED_PROVIDER`: provedor de preços. Padrão: `memory`. Use `twelvedata` para buscar cotações externas
- `TWELVEDATA_API_KEY`: chave da Twelve Data, obrigatória quando `PRICEFEED_PROVIDER=twelvedata`
- `TWELVEDATA_BASE_URL`: URL base da API da Twelve Data. Padrão: `https://api.twelvedata.com`

Exemplo com execução local sem banco relacional:

```bash
cd backend
PRICEFEED_PROVIDER=twelvedata \
TWELVEDATA_API_KEY=your-api-key \
go run ./cmd/api
```

Se `PRICEFEED_PROVIDER` não for definido, o backend continua usando o feed em memória para o MVP.
Se `DATABASE_URL` não for definido, o backend continua usando os repositórios em memória/arquivo do MVP.
No fluxo padrão com `docker compose`, `DATABASE_URL` já vem definido e a persistência relacional fica ativa.

## Endpoints

- `GET /healthz`
- `POST /auth/token`
- `GET /actions`
- `POST /actions`
- `GET /watchlist`
- `POST /watchlist`
- `DELETE /watchlist/{action_id}`
- `GET /alerts`
- `POST /alerts`
- `PATCH /alerts/{id}`
- `DELETE /alerts/{id}`
- `GET /devices`
- `POST /devices/register`
- `GET /prices`
- `PUT /prices`
- `POST /prices/check`

> `GET /actions` aceita `query` e faz busca exata por `name`.
> `POST /actions` cria ou reativa uma ação no catálogo usando `symbol`, `name` e `exchange`.
> O cadastro de alerta agora usa `action_id`, obtido em `GET /actions` ou `POST /actions`.
> Criar um alerta adiciona a ação na watchlist do usuário automaticamente, se necessário.

## Documentação Postman

A collection fica em `backend/docs/postman/`:

- `ideacoes-b3-alerts.postman_collection.json`
- `ideacoes-b3-alerts.postman_environment.json`

Fluxo recomendado no Postman:

1. `GET /actions`
2. `POST /actions`
3. `GET /watchlist`
4. `POST /watchlist`
5. `POST /devices/register`
6. `POST /alerts` com `action_id`
7. `PATCH /alerts/{id}`
8. `DELETE /alerts/{id}`

## Testes

Na raiz do projeto:

```bash
make test-backend
make test-backend-integration
```

O segundo comando sobe um Postgres temporário via Testcontainers e valida persistência real de alertas e devices.
