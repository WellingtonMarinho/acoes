import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../actions/domain/action.dart' as stock_action;
import '../../../core/network/api_client.dart';
import '../../alerts/presentation/alerts_controller.dart';
import '../../watchlist/domain/watchlist_item.dart';
import '../../watchlist/presentation/watchlist_controller.dart';

class WatchlistPageView extends ConsumerWidget {
  const WatchlistPageView({
    super.key,
    required this.onAddAction,
    required this.onCreateAlert,
  });

  final VoidCallback onAddAction;
  final ValueChanged<stock_action.MarketAction> onCreateAlert;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final watchlist = ref.watch(watchlistControllerProvider);

    return RefreshIndicator(
      onRefresh: () => ref.read(watchlistControllerProvider.notifier).refresh(),
      child: ListView(
        physics: const AlwaysScrollableScrollPhysics(),
        padding: const EdgeInsets.all(20),
        children: [
          Row(
            children: [
              Expanded(
                child: Text(
                  'Ações monitoradas',
                  style: Theme.of(context).textTheme.titleLarge?.copyWith(
                        fontWeight: FontWeight.w800,
                      ),
                ),
              ),
              FilledButton.icon(
                onPressed: onAddAction,
                icon: const Icon(Icons.add),
                label: const Text('Adicionar ação'),
              ),
            ],
          ),
          const SizedBox(height: 16),
          watchlist.when(
            data: (items) {
              if (items.isEmpty) {
                return _EmptyState(onAddAction: onAddAction);
              }

              return Column(
                children: [
                  for (final item in items) ...[
                    _WatchlistItemCard(
                      item: item,
                      onCreateAlert: onCreateAlert,
                      onRemove: () async {
                        await _removeAction(context, ref, item);
                      },
                    ),
                    const SizedBox(height: 12),
                  ],
                ],
              );
            },
            loading: () => const Padding(
              padding: EdgeInsets.symmetric(vertical: 40),
              child: Center(child: CircularProgressIndicator()),
            ),
            error: (error, stackTrace) => _ErrorState(
              message: error.toString(),
              onRetry: () =>
                  ref.read(watchlistControllerProvider.notifier).refresh(),
            ),
          ),
        ],
      ),
    );
  }

  Future<void> _removeAction(
    BuildContext context,
    WidgetRef ref,
    WatchlistItem item,
  ) async {
    final confirm = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Remover ação'),
        content: Text(
            'Remover ${item.symbol} da watchlist também exclui os alertas vinculados.'),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(false),
            child: const Text('Cancelar'),
          ),
          FilledButton(
            onPressed: () => Navigator.of(context).pop(true),
            child: const Text('Remover'),
          ),
        ],
      ),
    );
    if (confirm != true) {
      return;
    }

    try {
      await ref
          .read(activeWatchlistRepositoryProvider)
          .removeWatchlist(item.actionId);
      await ref.read(watchlistControllerProvider.notifier).refresh();
      await ref.read(alertsControllerProvider.notifier).refresh();
      if (context.mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('${item.symbol} removida da watchlist.')),
        );
      }
    } on ApiException catch (error) {
      if (context.mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text(ApiClient.friendlyMessageFor(error))),
        );
      }
    }
  }
}

class _WatchlistItemCard extends StatelessWidget {
  const _WatchlistItemCard({
    required this.item,
    required this.onCreateAlert,
    required this.onRemove,
  });

  final WatchlistItem item;
  final ValueChanged<stock_action.MarketAction> onCreateAlert;
  final Future<void> Function() onRemove;

  @override
  Widget build(BuildContext context) {
    final colors = Theme.of(context).colorScheme;
    final action = stock_action.MarketAction(
      id: item.actionId,
      symbol: item.symbol,
      name: item.name,
      exchange: item.exchange,
    );

    return Card(
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        item.symbol,
                        style: Theme.of(context).textTheme.titleLarge?.copyWith(
                              fontWeight: FontWeight.w800,
                            ),
                      ),
                      const SizedBox(height: 2),
                      Text(
                        item.name,
                        style: Theme.of(context).textTheme.bodyMedium,
                      ),
                      const SizedBox(height: 6),
                      Text(
                        item.exchange,
                        style:
                            Theme.of(context).textTheme.labelMedium?.copyWith(
                                  color: colors.onSurfaceVariant,
                                ),
                      ),
                    ],
                  ),
                ),
                if (item.openAlertsCount > 0)
                  Chip(
                    label: Text('${item.openAlertsCount} alerta(s)'),
                    backgroundColor: colors.secondaryContainer,
                  ),
              ],
            ),
            const SizedBox(height: 16),
            Row(
              children: [
                Expanded(
                  child: _Metric(
                    label: 'Preço atual',
                    value: item.currentPrice == null
                        ? 'Sem dado'
                        : 'R\$ ${item.currentPrice!.toStringAsFixed(2)}',
                  ),
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: _Metric(
                    label: 'Atualizado em',
                    value: item.lastPriceAt == null
                        ? 'Sem horário'
                        : _formatDateTime(item.lastPriceAt!),
                  ),
                ),
              ],
            ),
            const SizedBox(height: 16),
            Wrap(
              spacing: 8,
              runSpacing: 8,
              children: [
                FilledButton.tonalIcon(
                  onPressed: () => onCreateAlert(action),
                  icon: const Icon(Icons.add_alert_outlined),
                  label: const Text('Criar alerta'),
                ),
                OutlinedButton.icon(
                  onPressed: onRemove,
                  icon: const Icon(Icons.delete_outline),
                  label: const Text('Remover'),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }
}

class _Metric extends StatelessWidget {
  const _Metric({
    required this.label,
    required this.value,
  });

  final String label;
  final String value;

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          label,
          style: Theme.of(context).textTheme.labelMedium?.copyWith(
                color: Theme.of(context).colorScheme.onSurfaceVariant,
              ),
        ),
        const SizedBox(height: 4),
        Text(
          value,
          style: Theme.of(context).textTheme.bodyLarge?.copyWith(
                fontWeight: FontWeight.w700,
              ),
        ),
      ],
    );
  }
}

class _EmptyState extends StatelessWidget {
  const _EmptyState({required this.onAddAction});

  final VoidCallback onAddAction;

  @override
  Widget build(BuildContext context) {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(20),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              'Nenhuma ação monitorada ainda.',
              style: Theme.of(context).textTheme.titleMedium?.copyWith(
                    fontWeight: FontWeight.w800,
                  ),
            ),
            const SizedBox(height: 8),
            const Text(
                'Adicione uma ação para começar a acompanhar preços e alertas.'),
            const SizedBox(height: 16),
            FilledButton(
              onPressed: onAddAction,
              child: const Text('Adicionar ação'),
            ),
          ],
        ),
      ),
    );
  }
}

class _ErrorState extends StatelessWidget {
  const _ErrorState({
    required this.message,
    required this.onRetry,
  });

  final String message;
  final VoidCallback onRetry;

  @override
  Widget build(BuildContext context) {
    return Card(
      child: Padding(
        padding: const EdgeInsets.all(20),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              'Falha ao carregar watchlist',
              style: Theme.of(context).textTheme.titleMedium?.copyWith(
                    fontWeight: FontWeight.w800,
                  ),
            ),
            const SizedBox(height: 8),
            Text(message),
            const SizedBox(height: 16),
            FilledButton.tonal(
                onPressed: onRetry, child: const Text('Tentar novamente')),
          ],
        ),
      ),
    );
  }
}

String _formatDateTime(DateTime value) {
  final local = value.toLocal();
  final day = local.day.toString().padLeft(2, '0');
  final month = local.month.toString().padLeft(2, '0');
  final hour = local.hour.toString().padLeft(2, '0');
  final minute = local.minute.toString().padLeft(2, '0');
  return '$day/$month $hour:$minute';
}
