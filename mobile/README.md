# Mobile

App Flutter para:

- criar alertas
- listar watchlist
- configurar tokens de push

## Status

O projeto Flutter mínimo já está estruturado neste diretório.

O app já tem:

- home
- sessão provisória com token
- criação de alerta
- registro de device
- persistência local da sessão

Parte da navegação ainda usa dados demo enquanto a integração total com o backend é fechada.

## Estrutura

- `lib/app/`: bootstrap da aplicação, tema e navegação
- `lib/core/`: client HTTP, configuração e utilitários compartilhados
- `lib/features/`: funcionalidades por domínio
- `test/`: testes unitários e de widget

## Próximos passos

1. Consolidar a integração do client HTTP com o backend protegido.
2. Melhorar o pós-submit das telas principais.
3. Cobrir os fluxos com testes de widget e de repositório.
