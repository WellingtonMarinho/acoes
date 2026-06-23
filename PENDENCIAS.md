# Pendências Técnicas

Arquivo central para registrar dívidas técnicas e itens que não devem passar batido enquanto o MVP evolui.

## Acompanhamento do backend

- [ ] Substituir a emissão simples de JWT por um fluxo real de login/identidade confiável.
- [ ] Definir estratégia de push notifications reais em vez do notifier de log.
- [ ] Trocar o feed local de preços por uma fonte real de cotações da B3.
- [ ] Avaliar persistência em banco relacional quando o fluxo sair do MVP.
- [ ] Revisar prevenção de disparos duplicados no worker de monitoramento.
- [x] Cobrir o caminho principal do MVP com teste de integração ponta a ponta.

## Acompanhamento do mobile

- [ ] Refinar a experiência pós-submit nas telas de alerta e device.
- [ ] Adicionar watchlist e histórico visível na interface.
- [x] Cobrir com testes os fluxos principais do app Flutter.
- [x] Consolidar a integração do app com o backend protegido.

## Observações

- Quando uma pendência for resolvida, remova daqui e, se aplicável, atualize a documentação ou o roadmap.
- Se surgir uma nova dívida técnica relevante, registre primeiro neste arquivo para manter o rastreio centralizado.
