import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import '../../../core/network/api_client.dart';
import '../../alerts/domain/alert.dart';
import '../../alerts/presentation/alerts_controller.dart';
import '../../alerts/presentation/widgets/alert_status_chip.dart';
import '../../watchlist/presentation/watchlist_controller.dart';

enum _AlertSection { active, history }

class AlertsPageView extends ConsumerStatefulWidget {
  const AlertsPageView({
    super.key,
    required this.onEditAlert,
  });

  final ValueChanged<Alert> onEditAlert;

  @override
  ConsumerState<AlertsPageView> createState() => _AlertsPageViewState();
}

class _AlertsPageViewState extends ConsumerState<AlertsPageView> {
  _AlertSection _section = _AlertSection.active;

  @override
  Widget build(BuildContext context) {
    final alerts = ref.watch(alertsControllerProvider);
    final title = switch (_section) {
      _AlertSection.active => 'Alertas ativos',
      _AlertSection.history => 'Histórico de alertas',
    };

    return RefreshIndicator(
      onRefresh: () => ref.read(alertsControllerProvider.notifier).refresh(),
      child: ListView(
        physics: const AlwaysScrollableScrollPhysics(),
        padding: const EdgeInsets.all(20),
        children: [
          Text(
            title,
            style: Theme.of(context).textTheme.titleLarge?.copyWith(
                  fontWeight: FontWeight.w800,
                ),
          ),
          const SizedBox(height: 16),
          SegmentedButton<_AlertSection>(
            segments: const [
              ButtonSegment(
                value: _AlertSection.active,
                label: Text('Ativos'),
                icon: Icon(Icons.notifications_active_outlined),
              ),
              ButtonSegment(
                value: _AlertSection.history,
                label: Text('Histórico'),
                icon: Icon(Icons.history),
              ),
            ],
            selected: {_section},
            onSelectionChanged: (value) {
              setState(() {
                _section = value.first;
              });
            },
          ),
          const SizedBox(height: 16),
          alerts.when(
            data: (items) {
              final filtered = switch (_section) {
                _AlertSection.active => items
                    .where((alert) => alert.status == AlertStatus.open)
                    .toList(),
                _AlertSection.history => items
                    .where((alert) => alert.status == AlertStatus.triggered)
                    .toList(),
              };

              final emptyTitle = switch (_section) {
                _AlertSection.active => 'Nenhum alerta encontrado.',
                _AlertSection.history => 'Nenhum alerta disparado ainda.',
              };
              final emptyBody = switch (_section) {
                _AlertSection.active =>
                  'Crie alertas a partir de uma ação monitorada para acompanhar preços.',
                _AlertSection.history =>
                  'Os alertas disparados aparecem aqui com data e hora do evento.',
              };

              if (filtered.isEmpty) {
                return Card(
                  child: Padding(
                    padding: const EdgeInsets.all(20),
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          emptyTitle,
                          style: const TextStyle(
                              fontSize: 18, fontWeight: FontWeight.w800),
                        ),
                        const SizedBox(height: 8),
                        Text(emptyBody),
                      ],
                    ),
                  ),
                );
              }

              return Column(
                children: [
                  for (final alert in filtered) ...[
                    _AlertCard(
                      alert: alert,
                      onEdit: _section == _AlertSection.active
                          ? () => widget.onEditAlert(alert)
                          : null,
                      onDelete: () => _deleteAlert(context, ref, alert),
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
            error: (error, stackTrace) => Card(
              child: Padding(
                padding: const EdgeInsets.all(20),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    const Text(
                      'Falha ao carregar alertas',
                      style:
                          TextStyle(fontSize: 18, fontWeight: FontWeight.w800),
                    ),
                    const SizedBox(height: 8),
                    Text(error.toString()),
                    const SizedBox(height: 16),
                    FilledButton(
                      onPressed: () =>
                          ref.read(alertsControllerProvider.notifier).refresh(),
                      child: const Text('Tentar novamente'),
                    ),
                  ],
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }

  Future<void> _deleteAlert(
      BuildContext context, WidgetRef ref, Alert alert) async {
    final confirm = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Excluir alerta'),
        content: Text('Excluir definitivamente o alerta de ${alert.symbol}?'),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(false),
            child: const Text('Cancelar'),
          ),
          FilledButton(
            onPressed: () => Navigator.of(context).pop(true),
            child: const Text('Excluir'),
          ),
        ],
      ),
    );
    if (confirm != true) {
      return;
    }

    try {
      await ref
          .read(activeAlertsRepositoryProvider)
          .deleteAlert(alertId: alert.id);
      await ref.read(alertsControllerProvider.notifier).refresh();
      await ref.read(watchlistControllerProvider.notifier).refresh();
      if (context.mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('Alerta de ${alert.symbol} excluído.')),
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

class _AlertCard extends StatelessWidget {
  const _AlertCard({
    required this.alert,
    required this.onEdit,
    required this.onDelete,
  });

  final Alert alert;
  final VoidCallback? onEdit;
  final VoidCallback onDelete;

  @override
  Widget build(BuildContext context) {
    final directionLabel = switch (alert.direction) {
      AlertDirection.above => 'acima',
      AlertDirection.below => 'abaixo',
    };

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
                        alert.actionName.isEmpty
                            ? alert.symbol
                            : alert.actionName,
                        style:
                            Theme.of(context).textTheme.titleMedium?.copyWith(
                                  fontWeight: FontWeight.w800,
                                ),
                      ),
                      const SizedBox(height: 4),
                      Text(
                          '${alert.symbol} · R\$ ${alert.targetPrice.toStringAsFixed(2)} · $directionLabel'),
                      if (alert.triggeredAt != null) ...[
                        const SizedBox(height: 4),
                        Text(
                          'Disparado em ${_formatTimestamp(alert.triggeredAt!)}',
                          style:
                              Theme.of(context).textTheme.bodySmall?.copyWith(
                                    color: Theme.of(context)
                                        .colorScheme
                                        .onSurfaceVariant,
                                  ),
                        ),
                      ],
                    ],
                  ),
                ),
                AlertStatusChip(status: alert.status),
              ],
            ),
            const SizedBox(height: 16),
            Row(
              children: [
                if (onEdit != null) ...[
                  FilledButton.tonalIcon(
                    onPressed: onEdit,
                    icon: const Icon(Icons.edit_outlined),
                    label: const Text('Editar'),
                  ),
                  const SizedBox(width: 8),
                ],
                OutlinedButton.icon(
                  onPressed: onDelete,
                  icon: const Icon(Icons.delete_outline),
                  label: const Text('Excluir'),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }

  String _formatTimestamp(DateTime timestamp) {
    final local = timestamp.toLocal();
    final day = local.day.toString().padLeft(2, '0');
    final month = local.month.toString().padLeft(2, '0');
    final year = local.year.toString();
    final hour = local.hour.toString().padLeft(2, '0');
    final minute = local.minute.toString().padLeft(2, '0');
    return '$day/$month/$year às $hour:$minute';
  }
}
