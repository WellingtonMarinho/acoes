# Backend

API em Go para alertas de ações da B3.

## Como subir localmente

Na raiz do projeto, rode:

```bash
cd backend
go run ./cmd/api
```

Por padrão, o servidor sobe em `http://localhost:8080`.

## Variáveis de ambiente

- `JWT_SECRET`: segredo usado para assinar e validar os JWTs. Obrigatório
- `HTTP_ADDR`: endereço HTTP do servidor. Padrão: `:8080`
- `MONITOR_INTERVAL_SECONDS`: intervalo do worker de monitoramento. Padrão: `10`
- `ALERTS_STORE_PATH`: caminho do arquivo de persistência dos alertas
- `DEVICES_STORE_PATH`: caminho do arquivo de persistência dos devices

Exemplo com persistência em arquivo:

```bash
cd backend
JWT_SECRET=dev-secret \
ALERTS_STORE_PATH=./data/alerts.json \
DEVICES_STORE_PATH=./data/devices.json \
go run ./cmd/api
```

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
