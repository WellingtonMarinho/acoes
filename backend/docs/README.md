# Backend Docs

Arquivos para validar a API do backend no Postman.

## Arquivos

- `postman/ideacoes-b3-alerts.postman_collection.json`
- `postman/ideacoes-b3-alerts.postman_environment.json`

## Como usar

1. Importa a collection no Postman.
2. Importa o environment local.
3. Seleciona o environment `Ideacoes Local`.
4. Executa `POST /auth/token` e guarda o `accessToken` no environment.
5. Registra o device token via `POST /devices/register`.
6. Cria alertas via `POST /alerts`.
7. Atualiza o preço via `PUT /prices`.
8. Aguarda o worker ou usa `POST /prices/check` para validação manual.

> As rotas `GET /alerts`, `GET /devices`, `POST /alerts` e `POST /devices/register` exigem `Authorization: Bearer <token>`.

## Base URL

Por padrão o backend roda em `http://localhost:8080`.
