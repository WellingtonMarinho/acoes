# Plano de Acao - Catalogo de Acoes e Watchlist

## Objetivo

Reorganizar o produto para um modelo coerente com a B3:

1. O sistema possui um catalogo de acoes preexistentes.
2. O usuario nao cadastra acoes manualmente.
3. O usuario escolhe quais acoes quer monitorar.
4. O usuario define os alertas de preco alvo para as acoes monitoradas.
5. O sistema se autoalimenta com as acoes disponiveis e pode evoluir para sync automatico.

## Premissas de Produto

- O universo inicial inclui apenas acoes ON, PN e units.
- BDRs, ETFs e FIIs nao fazem parte do MVP.
- Somente acoes ativas devem aparecer no catalogo.
- O MVP pode iniciar com cerca de 20 acoes fixas.
- A arquitetura deve facilitar crescimento para uma lista maior e sincronizacao automatica.
- Busca deve funcionar por simbolo e por nome.
- Quando o usuario tentar adicionar uma acao inexistente no catalogo, o backend deve tentar buscar e autocadastrar essa acao.
- Criar alerta pode aproveitar a mesma camada de auto ingestao.
- A atualizacao do catalogo sera inicialmente via rotina offline periodica.

## Direcao Tecnica

### Catalogo de Acoes

- `actions` passa a ser um catalogo global somente leitura para o usuario final.
- O usuario nao cria acoes via app.
- O backend e responsavel por manter esse catalogo atualizado.
- O catalogo deve suportar:
  - listagem
  - busca por simbolo
  - busca por nome
  - upsert interno para ingestao

### Watchlist

- `watchlist` representa as acoes que cada usuario deseja monitorar.
- O usuario adiciona e remove acoes da propria watchlist.
- A watchlist deve continuar sendo o ponto de entrada para o monitoramento do preco.
- A remocao de uma acao da watchlist deve limpar os alertas vinculados quando a regra de negocio exigir cascata.

### Alerts

- Alertas continuam sendo criados a partir de uma acao monitorada.
- O fluxo de alerta deve garantir que a acao exista no catalogo.
- O usuario deve poder editar e excluir alertas.
- Exclusao deve ser fisica.

## Fluxo Esperado

1. O backend sobe com o catalogo inicial carregado.
2. O usuario busca uma acao por simbolo ou nome.
3. Se a acao existir, ela pode ser adicionada a watchlist.
4. Se a acao nao existir, o backend tenta ingestao/autocadastro antes de concluir a operacao.
5. O usuario cria um alerta para a acao monitorada.
6. O sistema reaproveita a mesma base de catalogo para novas acoes que surgirem depois.

## Backend

### Entregas

- Remover o fluxo publico de cadastro manual de acao.
- Manter `actions` como catalogo global.
- Implementar busca por simbolo e por nome.
- Preparar upsert interno para ingestao automatica.
- Criar estrutura de seed inicial com cerca de 20 acoes.
- Preparar rotina offline/cron para atualizacao do catalogo.
- Manter watchlist por usuario.
- Garantir que criar alerta possa autocadastrar a acao se necessario.

### Contratos

- `GET /actions`
- `GET /watchlist`
- `POST /watchlist`
- `DELETE /watchlist/{action_id}`
- `GET /alerts`
- `POST /alerts`
- `PATCH /alerts/{id}`
- `DELETE /alerts/{id}`

### Regras

- Buscar por simbolo e nome.
- Exibir apenas acoes ativas.
- Nao expor cadastro manual de acoes ao usuario.
- Evitar duplicidade de watchlist.
- Evitar alertas orfaos.

## Mobile

### Entregas

- Remover a ideia de "cadastrar acao" manual.
- Exibir busca e selecao do catalogo.
- Permitir adicionar acao existente a watchlist.
- Permitir criar alerta a partir da acao monitorada.
- Exibir feedback claro em sucesso, erro e loading.
- Manter o app leve, escaneavel e com foco em monitoramento.

### UX

- A tela principal deve priorizar watchlist e alertas.
- O usuario deve perceber que o sistema ja conhece o catalogo de acoes.
- O fluxo de adicionar deve parecer "selecionar e monitorar", nao "cadastrar um ativo".

## Testes

- Cobrir o catalogo e a watchlist no backend.
- Cobrir ingestao/autocadastro com teste automatizado.
- Cobrir os fluxos principais do mobile com widget tests.
- Validar que o app reage bem quando a acao ainda nao existe no catalogo.

## Proximas Etapas

1. Remover o cadastro manual de acao do mobile.
2. Ajustar o backend para catalogo orientado a ingestao.
3. Criar seed inicial com as acoes essenciais.
4. Implementar a busca por simbolo e nome.
5. Preparar a rotina offline de sincronizacao do catalogo.
