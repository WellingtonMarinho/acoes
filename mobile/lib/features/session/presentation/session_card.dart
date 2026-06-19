import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

import 'session_controller.dart';

class SessionCard extends ConsumerStatefulWidget {
  const SessionCard({super.key});

  @override
  ConsumerState<SessionCard> createState() => _SessionCardState();
}

class _SessionCardState extends ConsumerState<SessionCard> {
  final _controller = TextEditingController(text: 'user-001');

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  Future<void> _signIn() async {
    if (_controller.text.trim().isEmpty) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Informe o user_id.')),
      );
      return;
    }

    await ref.read(sessionControllerProvider.notifier).signIn(
          userId: _controller.text.trim(),
        );
  }

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
              session == null ? 'Entrar no app' : 'Sessão ativa',
              style: Theme.of(context).textTheme.titleMedium?.copyWith(
                    fontWeight: FontWeight.w700,
                  ),
            ),
            const SizedBox(height: 8),
            Text(
              session == null
                  ? 'Use um user_id provisório para emitir o JWT do MVP.'
                  : 'Usuário ${session.userId} autenticado com token do MVP.',
            ),
            const SizedBox(height: 12),
            if (session == null) ...[
              TextField(
                controller: _controller,
                decoration: const InputDecoration(
                  labelText: 'User ID',
                  hintText: 'user-001',
                ),
              ),
              const SizedBox(height: 12),
              FilledButton(
                onPressed: _signIn,
                child: const Text('Emitir token'),
              ),
            ] else ...[
              OutlinedButton(
                onPressed: () => ref.read(sessionControllerProvider.notifier).signOut(),
                child: const Text('Sair'),
              ),
            ],
          ],
        ),
      ),
    );
  }
}
