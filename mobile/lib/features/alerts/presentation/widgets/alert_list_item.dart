import 'package:flutter/material.dart';

import '../../domain/alert.dart';
import 'alert_status_chip.dart';

class AlertListItem extends StatelessWidget {
  const AlertListItem({
    super.key,
    required this.alert,
  });

  final Alert alert;

  @override
  Widget build(BuildContext context) {
    final directionLabel = switch (alert.direction) {
      AlertDirection.above => 'acima',
      AlertDirection.below => 'abaixo',
    };

    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(18),
        border: Border.all(color: Theme.of(context).colorScheme.outlineVariant),
      ),
      child: Row(
        children: [
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  alert.symbol,
                  style: Theme.of(context).textTheme.titleMedium?.copyWith(
                        fontWeight: FontWeight.w700,
                      ),
                ),
                const SizedBox(height: 4),
                Text('R\$ ${alert.targetPrice.toStringAsFixed(2)} · $directionLabel'),
              ],
            ),
          ),
          AlertStatusChip(status: alert.status),
        ],
      ),
    );
  }
}
