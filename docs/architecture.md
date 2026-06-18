    # Arquitetura MVP

## Fluxo

1. O app registra o `device_token` do usuário no backend.
2. O usuário cria um alerta para um ativo e um preço alvo.
3. Um endpoint atualiza o preço corrente do ativo no feed local.
4. Um worker consulta os preços e compara com os alertas abertos.
5. Quando o alvo é atingido, o alerta muda para `triggered`.
6. O backend dispara notificação push para o dispositivo registrado.

## Componentes

- App Flutter
- API Go
- Repositório de alertas
- Serviço de cotação
- Worker de avaliação
- Serviço de push
- Registro de devices por usuário

## Decisões intencionais

- Manter o domínio de alertas isolado do transporte HTTP.
- Começar com armazenamento em memória para acelerar o MVP.
- Usar interfaces para plugar fonte de dados e push depois sem reescrever o core.
- Suportar persistência em arquivo no MVP para não perder alertas ao reiniciar a API.
- Suportar persistência em arquivo para registros de devices.
