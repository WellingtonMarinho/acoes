# Roadmap do Produto

Este documento consolida o caminho do produto Ideações, do que já foi feito até o que ainda vamos construir.

## Visão

O produto é um app para monitorar ações da B3 e disparar alertas de preço por push notification quando um ativo atingir um valor alvo, para cima ou para baixo.

O stack inicial é:

- Mobile em Flutter
- Backend em Go
- Push com FCM/APNs
- Persistência local no MVP, com evolução natural para banco relacional

## Estado atual

### Já entregue

- Estrutura inicial do projeto com `backend/`, `mobile/` e `docs/`
- API Go base com healthcheck e endpoints principais
- Domínio de alertas isolado em regras de negócio
- Criação e listagem de alertas
- Disparo de alerta quando o preço cruza o alvo
- Repositório em memória para alertas
- Persistência em arquivo para alertas via `ALERTS_STORE_PATH`
- Registro de device token por usuário
- Persistência em arquivo para devices via `DEVICES_STORE_PATH`
- Feed local de preços para simulação no MVP
- Worker de monitoramento em loop com intervalo configurável
- Collection exportável do Postman para validar a API
- Testes automatizados para alertas, devices e persistência

### O que isso já permite hoje

- Registrar o token do dispositivo
- Criar um alerta de alta ou baixa
- Simular o preço de um ativo
- Disparar o alerta automaticamente quando o preço cruza o alvo
- Validar manualmente pelo Postman
- Manter alertas e devices entre reinícios, se a persistência em arquivo estiver ativa

## Roadmap por fases

### Fase 1. MVP técnico

Objetivo: ter o núcleo funcional do produto, testável localmente.

Entregas:

- API com alertas
- Persistência básica
- Registro de device token
- Monitoramento de preços
- Collection do Postman
- Documentação mínima do backend

Status:

- Concluída

### Fase 2. Integração com o app mobile

Objetivo: transformar a API em uma experiência real de uso.

Entregas:

- Projeto Flutter inicial
- Tela de login ou identificação básica do usuário
- Tela de cadastro de device token
- Tela de watchlist
- Tela de criação de alertas
- Tela de histórico de alertas disparados
- Consumo da API do backend

Status:

- Ainda não iniciado

### Fase 3. Push notifications de verdade

Objetivo: sair do log e enviar notificações reais para o celular.

Entregas:

- Integração com FCM no Android
- Integração com APNs no iPhone
- Interface de notifier desacoplada do domínio
- Envio assíncrono ou via fila/outbox
- Registro de tentativas e falhas de entrega

Status:

- Estrutura preparada, integração real ainda pendente

### Fase 4. Fonte real de cotações

Objetivo: substituir o feed local por dados reais de mercado.

Entregas:

- Integração com uma fonte confiável de cotações da B3
- Normalização do formato de preço por ativo
- Controle de taxa de requisições
- Cache para reduzir custo e latência
- Estratégia de fallback se a fonte ficar indisponível

Status:

- Ainda não iniciado

### Fase 5. Robustez operacional

Objetivo: deixar o produto pronto para uso contínuo.

Entregas:

- Banco relacional para alertas, devices e histórico
- Migrações versionadas
- Observabilidade com logs estruturados e métricas
- Tratamento de concorrência para evitar disparos duplicados
- Reprocessamento seguro de alertas
- Ambientes separados por configuração

Status:

- Ainda não iniciado

### Fase 6. Produto e experiência

Objetivo: melhorar retenção e usabilidade.

Entregas:

- Categorias de alertas
- Alertas recorrentes
- Lista de favoritos/watchlist
- Filtros por ativo e status
- Histórico de preço e alertas
- Melhor UX para criar alertas em poucos toques
- Preferências por usuário

Status:

- Ainda não iniciado

## Ordem sugerida de execução

1. Implementar o app Flutter mínimo.
2. Integrar o backend com push notifications reais.
3. Conectar uma fonte real de cotações.
4. Migrar persistência para banco relacional.
5. Adicionar observabilidade e prevenção de duplicidade.

## Critérios de pronto por etapa

### Backend mínimo pronto

- Cria alertas
- Lista alertas
- Registra devices
- Monitora preços
- Dispara notificação

### App mínimo pronto

- Conecta na API
- Registra token
- Permite criar alertas
- Mostra alertas e status

### Produto pronto para beta

- Push real funcionando
- Preço real funcionando
- Persistência confiável
- Fluxo de criação simples
- Sem perda de dados em restart

## Decisões já tomadas

- Flutter no mobile
- Go no backend
- Domínio separado do transporte HTTP
- Modo local para simulação e teste
- Persistência em arquivo no MVP para acelerar validação

## Próximos passos recomendados

1. Iniciar o app Flutter com navegação básica.
2. Definir a integração de push real.
3. Escolher a primeira fonte de preços da B3.

