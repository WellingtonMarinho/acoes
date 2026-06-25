import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'session_controller.dart';

class SessionCard extends ConsumerStatefulWidget {
  const SessionCard({super.key});

  @override
  ConsumerState<SessionCard> createState() => _SessionCardState();
}

class _SessionCardState extends ConsumerState<SessionCard> {
  @override
  Widget build(BuildContext context) {
    final session = ref.watch(sessionControllerProvider);

    return Card(
      child: Padding(
        padding: const EdgeInsets.all(20),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              'Sessão automática',
              style: Theme.of(context).textTheme.titleMedium?.copyWith(
                    fontWeight: FontWeight.w700,
                  ),
            ),
            const SizedBox(height: 8),
            Text(
              session == null
                  ? 'Não foi possível iniciar a sessão automaticamente.'
                  : 'Usuário ${session.userId} autenticado com token do MVP.',
            ),
          ],
        ),
      ),
    );
  }
}
