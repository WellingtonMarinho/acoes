import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';

import '../../../app/theme_controller.dart';
import '../../session/presentation/session_card.dart';

class SettingsPageView extends ConsumerWidget {
  const SettingsPageView({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final themeModeAsync = ref.watch(themeModeControllerProvider);

    return ListView(
      padding: const EdgeInsets.all(20),
      children: [
        const SessionCard(),
        const SizedBox(height: 20),
        Card(
          child: Padding(
            padding: const EdgeInsets.all(20),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  'Tema',
                  style: Theme.of(context).textTheme.titleMedium?.copyWith(
                        fontWeight: FontWeight.w800,
                      ),
                ),
                const SizedBox(height: 8),
                const Text('Escolha entre sistema, claro ou escuro.'),
                const SizedBox(height: 16),
                themeModeAsync.when(
                  data: (mode) => SegmentedButton<ThemeMode>(
                    segments: const [
                      ButtonSegment(
                          value: ThemeMode.system, label: Text('Sistema')),
                      ButtonSegment(
                          value: ThemeMode.light, label: Text('Claro')),
                      ButtonSegment(
                          value: ThemeMode.dark, label: Text('Escuro')),
                    ],
                    selected: {mode},
                    onSelectionChanged: (value) {
                      ref
                          .read(themeModeControllerProvider.notifier)
                          .setThemeMode(value.first);
                    },
                  ),
                  loading: () => const LinearProgressIndicator(),
                  error: (error, stackTrace) => Text(error.toString()),
                ),
              ],
            ),
          ),
        ),
        const SizedBox(height: 20),
        Card(
          child: ListTile(
            title: const Text('Registrar device'),
            subtitle: const Text('Configurar token de push do usuário ativo'),
            trailing: const Icon(Icons.chevron_right),
            onTap: () => context.push('/devices/register'),
          ),
        ),
      ],
    );
  }
}
