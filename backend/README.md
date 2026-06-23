# Backend

API em Go para alertas de ações da B3.

## Como subir localmente

Na raiz do projeto:

```bash
make run-backend
```

Isso sobe o backend em `http://localhost:8080` e um Postgres em `localhost:5432`.
Quando `DATABASE_URL` estiver definido, o backend usa o Postgres para alertas e devices.

## Variáveis de ambiente

- `JWT_SECRET`: segredo usado para assinar e validar os JWTs. Obrigatório
- `HTTP_ADDR`: endereço HTTP do servidor. Padrão: `:8080`
- `MONITOR_INTERVAL_SECONDS`: intervalo do worker de monitoramento. Padrão: `10`
- `ALERTS_STORE_PATH`: caminho do arquivo de persistência dos alertas
- `DEVICES_STORE_PATH`: caminho do arquivo de persistência dos devices
- `DATABASE_URL`: string de conexão do Postgres. Quando definida, o backend usa persistência relacional
- `PRICEFEED_PROVIDER`: provedor de preços. Padrão: `memory`. Use `twelvedata` para buscar cotações externas
- `TWELVEDATA_API_KEY`: chave da Twelve Data, obrigatória quando `PRICEFEED_PROVIDER=twelvedata`
- `TWELVEDATA_BASE_URL`: URL base da API da Twelve Data. Padrão: `https://api.twelvedata.com`

Exemplo com persistência em arquivo:

```bash
cd backend
JWT_SECRET=dev-secret \
ALERTS_STORE_PATH=./data/alerts.json \
DEVICES_STORE_PATH=./data/devices.json \
PRICEFEED_PROVIDER=twelvedata \
TWELVEDATA_API_KEY=your-api-key \
go run ./cmd/api
```

Se `PRICEFEED_PROVIDER` não for definido, o backend continua usando o feed em memória para o MVP.
Se `DATABASE_URL` não for definido, o backend continua usando os repositórios em memória/arquivo do MVP.

## Endpoints

- `GET /healthz`
- `POST /auth/token`
- `GET /alerts`
- `POST /alerts`
- `GET /devices`
- `POST /devices/register`
- `GET /prices`
- `PUT /prices`
- `POST /prices/check`

> Nota: `GET /alerts`, `GET /devices`, `POST /alerts` e `POST /devices/register` agora exigem `Authorization: Bearer <token>`.
> O token é emitido em `POST /auth/token` usando `user_id` apenas nessa etapa inicial do MVP.

## Documentação Postman

A collection fica em `backend/docs/postman/`:

- `ideacoes-b3-alerts.postman_collection.json`
- `ideacoes-b3-alerts.postman_environment.json`

## Testes

Na raiz do projeto:

```bash
make test-backend
```
