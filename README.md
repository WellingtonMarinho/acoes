# Ideacoes

MVP de alertas para ações da B3 com backend em Go e app mobile em Flutter.

## Objetivo inicial

- Cadastrar alertas de preço por ativo.
- Receber notificações quando o preço cruzar o alvo para cima ou para baixo.
- Preparar a base para integrar um feed real de cotações e push notifications.

## Estrutura

- `backend/`: API e regras de negócio em Go.
- `mobile/`: app Flutter em construção.
- `docs/`: decisões de produto e arquitetura.

## Backend

O backend inicial já expõe:

- `GET /healthz`
- `POST /auth/token`
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

1. Emita um token em `POST /auth/token`.
2. Registre o device token em `POST /devices/register`.
3. Crie um alerta em `POST /alerts`.
4. Atualize um preço em `PUT /prices`.
5. Aguarde o worker ou force a checagem com `POST /prices/check`.

> `GET /alerts`, `GET /devices`, `POST /alerts` e `POST /devices/register` exigem `Authorization: Bearer <token>`.

## Mobile

O app Flutter já está sendo estruturado com:

- home
- tela de criação de alerta
- tela de registro de device
- sessão provisória com persistência local

O próximo passo é consolidar a integração do mobile com o backend protegido.

## CI/CD

O repositório tem pipelines separados por stack em `.github/workflows/`.

### Backend Go

- Executa em `push` para `main` e em `pull_request`
- Roda `go mod download`
- Roda `go test ./... -race -coverprofile=coverage.out`
- Roda `go vet ./...`
- Roda `golangci-lint`
- Roda `gosec`
- Publica `backend/coverage.out` como artifact

Comandos locais equivalentes:

```bash
cd backend
go mod download
go test ./... -race -coverprofile=coverage.out
go vet ./...
golangci-lint run
gosec ./...
```

### Mobile Flutter

- Executa em `push` para `main` e em `pull_request`
- Roda `flutter pub get`
- Roda `dart format --set-exit-if-changed .`
- Roda `flutter analyze`
- Roda `flutter test --coverage`
- Publica `mobile/coverage/lcov.info` como artifact

Comandos locais equivalentes:

```bash
cd mobile
flutter pub get
dart format --set-exit-if-changed .
flutter analyze
flutter test --coverage
```

### Cobertura

- O backend gera `coverage.out`
- O mobile gera `coverage/lcov.info`
- Esses arquivos ficam disponíveis nos artifacts do GitHub Actions para inspeção ou integração futura com serviços externos

## Próximo passo

Ligar uma fonte real de cotações da B3 e trocar o notifier de log por push via FCM/APNs.
# acoes
