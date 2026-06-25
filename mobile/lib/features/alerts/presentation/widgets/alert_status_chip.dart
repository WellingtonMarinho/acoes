import 'package:flutter/material.dart';

import '../../domain/alert.dart';

class AlertStatusChip extends StatelessWidget {
  const AlertStatusChip({
    super.key,
    required this.status,
  });

  final AlertStatus status;

  @override
  Widget build(BuildContext context) {
    final colors = Theme.of(context).colorScheme;
    final label = switch (status) {
      AlertStatus.open => 'ativo',
      AlertStatus.triggered => 'disparado',
    };
    final background = switch (status) {
      AlertStatus.open => colors.secondaryContainer,
      AlertStatus.triggered => colors.primaryContainer,
    };

    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 6),
      decoration: BoxDecoration(
        color: background,
        borderRadius: BorderRadius.circular(999),
      ),
      child: Text(
        label,
        style: Theme.of(context).textTheme.labelSmall?.copyWith(
              fontWeight: FontWeight.w700,
            ),
      ),
    );
  }
}
