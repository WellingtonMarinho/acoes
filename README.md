# Ideacoes

MVP de alertas para ações da B3 com backend em Go e app mobile em Flutter.

## Objetivo inicial

- Cadastrar alertas de preço por ativo.
- Receber notificações quando o preço cruzar o alvo para cima ou para baixo.
- Preparar a base para integrar um feed real de cotações e push notifications.

## Estrutura

- `backend/`: API e regras de negócio em Go.
- `mobile/`: app Flutter que virá na próxima etapa.
- `docs/`: decisões de produto e arquitetura.

## Backend

O backend inicial já expõe:

- `GET /healthz`
- `POST /alerts`
- `GET /alerts`
- `GET /devices`
- `POST /devices/register`
- `GET /prices`
- `PUT /prices`
- `POST /prices/check`

Por padrão ele usa armazenamento em memória. Para persistir os alertas, defina:

```bash
ALERTS_STORE_PATH=./data/alerts.json
```

Para persistir os devices, defina:

```bash
DEVICES_STORE_PATH=./data/devices.json
```

O worker de monitoramento usa `MONITOR_INTERVAL_SECONDS` e, por padrão, roda a cada 10 segundos.

Fluxo sugerido para teste:

1. Registre o device token em `POST /devices/register`.
2. Crie um alerta em `POST /alerts`.
3. Atualize um preço em `PUT /prices`.
4. Aguarde o worker ou force a checagem com `POST /prices/check`.

## Próximo passo

Ligar uma fonte real de cotações da B3 e trocar o notifier de log por push via FCM/APNs.
# acoes
