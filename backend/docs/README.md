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
5. Consulta `GET /actions` para pegar o `action_id`.
6. Usa `GET /watchlist` e `POST /watchlist` para montar a lista de monitoramento.
7. Registra o device token via `POST /devices/register`.
8. Cria alertas via `POST /alerts` com `action_id`.
9. Atualiza alertas via `PATCH /alerts/{id}` ou remove via `DELETE /alerts/{id}`.
10. Atualiza o preço via `PUT /prices`.
11. Aguarda o worker ou usa `POST /prices/check` para validação manual.

> As rotas `GET /watchlist`, `POST /watchlist`, `DELETE /watchlist/{action_id}`, `GET /alerts`, `GET /devices`, `POST /alerts`, `PATCH /alerts/{id}`, `DELETE /alerts/{id}` e `POST /devices/register` exigem `Authorization: Bearer <token>`.

## Base URL

Por padrão o backend roda em `http://localhost:8080`.
